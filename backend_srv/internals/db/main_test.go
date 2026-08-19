package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var SetupTestdbpg *pgxpool.Pool

// NOTE : Hook for starting all the test, and shared only in this package
func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	SetupTestDB1()

	os.Exit(m.Run())
}

func SetupTestDB1() {
	log.Printf("TestMain pkg db")

	ctxBg := context.Background()
	// variable dipake nanti
	conn, err := NewConn(ctxBg, os.Getenv("PG_CONNSTRING"))
	if err != nil {
		log.Printf("error init DB ====>: %s\n", err.Error())
	}

	log.Println("TEST PRINT HERE!")

	SetupTestdbpg = conn.DBPool

}
