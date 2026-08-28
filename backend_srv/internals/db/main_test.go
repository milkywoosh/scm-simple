package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"scm-simple-luke.com/dir/internals/utils"
)

var SetupTestdbpg *pgxpool.Pool
var config utils.Config

// NOTE : Hook for starting all the test, and shared only in this package
func TestMain(m *testing.M) {
	var err error
	config, err = utils.LoadConfig("../../../")
	if err != nil {
		os.Exit(m.Run())
	}
	SetupTestDB1()

	os.Exit(m.Run())
}

func SetupTestDB1() {
	log.Printf("TestMain pkg db")

	ctxBg := context.Background()
	// variable dipake nanti
	conn, err := NewConn(ctxBg, config.DBSource)
	if err != nil {
		log.Printf("error init DB ====>: %s\n", err.Error())
	}

	log.Println("TEST PRINT HERE!")

	SetupTestdbpg = conn.DBPool

}
