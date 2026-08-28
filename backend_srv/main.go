package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	_ "scm-simple-luke.com/dir/docs"
	"scm-simple-luke.com/dir/internals/api"
	"scm-simple-luke.com/dir/internals/db"
	"scm-simple-luke.com/dir/internals/token"
	"scm-simple-luke.com/dir/internals/utils"

	"scm-simple-luke.com/dir/internals/services"
)

/* swag init --parseDependency --parseDependencyLevel 1 */

// @title           scm-simple API
// @version         1.0
// @description     This is my service API documentation for scm-simple.
// @host            103.235.75.94:8000
// @BasePath        /api/v1

// @contact.name   API Support
// @contact.email  archer.lukman@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1
// @securityDefinitions.basic  BasicAuth

var Conn_ *db.Conn
var MainConn *sql.DB

// signal yg dikirim OS, hasil trigger dari inputan user/human
var interrupsSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func runRestAPIServer(
	ctx context.Context,
	config utils.Config,
	waitGroup *errgroup.Group,
	service *services.Services,
	taskDistributor any, // worker.TaskDistributor
) {

	newMaker, err := token.NewJWTMaker("7c2f6152b065929f57da4df44bc704d8")
	if err != nil {
		log.Fatalf("err info token maker: %s", err.Error())
	}
	serverREST, err := api.NewServer(config, service, newMaker, taskDistributor)
	if err != nil {
		log.Fatalf("err info server: %v", err)
	}

	// Start the server immediately

	waitGroup.Go(func() error {
		log.Printf("starting HTTP server on %s", config.Port)

		if err := serverREST.Start(config.Port); err != nil {
			log.Printf("HTTP server error: %v", err)
			return err
		}

		// note: this line will never reached, because  this actually block (right behaviour) => serverREST.Start(":8081")
		log.Printf("[this log never block] success starting HTTP server")
		return nil
	})

	waitGroup.Go(func() error {
		log.Printf("blocking cancelation signal from OS process...")
		// Wait for cancellation signal (CTRL+C, SIGTERM, etc.)
		<-ctx.Done()
		log.Printf("after receive cancelation signal from OS process...")

		log.Println("graceful shutdown HTTP server...")

		// give sometime to active connection before server really close
		log.Printf("give sometime to active connection before server really close")
		shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		if err := serverREST.Shutdown(shutdownCtx); err != nil {
			log.Printf("forced shutdown: %v", err)
			return err
		}

		log.Println("HTTP server stopped cleanly")
		return nil
	})

}

func main() {

	config, err := utils.LoadConfig("../")
	if err != nil {
		log.Printf("err load config: %s", err.Error())
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), interrupsSignal...)
	defer stop()

	// variable dipake nanti
	Conn_, err = db.NewConn(ctx, config.DBSource)
	// log.Printf("log err: %s", err.Error())
	if err != nil {
		log.Printf("error init DB ====>: %s\n", err.Error())
	}

	fmt.Println("go pg auth test")

	dbInstance := db.NewPgInstance(Conn_.DBPool)
	newServices := services.NewServices(dbInstance)

	waitGroup, ctxWg := errgroup.WithContext(ctx)

	// process2, task will be Enqueued if CALLED via http req
	runRestAPIServer(ctxWg, config, waitGroup, newServices, nil)

	// process3
	// ================ get err after wait several process
	err = waitGroup.Wait()
	if err != nil {
		log.Fatalf("last tail err info: %v", err)
	}

}
