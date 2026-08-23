package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	_ "scm-simple-luke.com/dir/docs"
	"scm-simple-luke.com/dir/internals/api"
	"scm-simple-luke.com/dir/internals/db"

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
	waitGroup *errgroup.Group,
	service *services.Services,
	taskDistributor any, // worker.TaskDistributor
) {
	serverREST, err := api.NewServer(service, taskDistributor)
	if err != nil {
		log.Fatalf("err info server: %v", err)
	}

	// Start the server immediately

	waitGroup.Go(func() error {
		log.Printf("starting HTTP server on %s", os.Getenv("PORT"))

		if err := serverREST.Start(os.Getenv("PORT")); err != nil {
			log.Printf("HTTP server error: %v", err)
			return err
		}

		// note: this line will never reached, because  this actually block (right behaviour) => serverREST.Start(":8081")
		log.Printf("[this log never block] success starting HTTP server")
		return nil
	})

	waitGroup.Go(func() error {
		// Wait for cancellation signal (CTRL+C, SIGTERM, etc.)
		<-ctx.Done()

		log.Println("graceful shutdown HTTP server...")

		if err := serverREST.Shutdown(ctx); err != nil {
			log.Printf("forced shutdown: %v", err)
			return err
		}

		log.Println("HTTP server stopped cleanly")
		return nil
	})
}

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("err load .env: %v ==> **only used in local dev, ignore on production", err)
	}
	fmt.Printf("%s\n", "cuss 1 scm")

	ctx, stop := signal.NotifyContext(context.Background(), interrupsSignal...)
	defer stop()

	// variable dipake nanti
	Conn_, err = db.NewConn(ctx, os.Getenv("PG_CONNSTRING"))
	// log.Printf("log err: %s", err.Error())
	if err != nil {
		log.Printf("error init DB ====>: %s\n", err.Error())
	}

	fmt.Println("go pg auth test")

	dbInstance := db.NewPgInstance(Conn_.DBPool)
	newServices := services.NewServices(dbInstance)

	waitGroup, ctxWg := errgroup.WithContext(ctx)

	// process2, task will be Enqueued if CALLED via http req
	runRestAPIServer(ctxWg, waitGroup, newServices, nil)

	// process3
	// ================ get err after wait several process
	err = waitGroup.Wait()
	if err != nil {
		log.Fatalf("last tail err info: %v", err)
	}

}
