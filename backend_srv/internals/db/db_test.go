package db

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestDBPositiveCase1(t *testing.T) {

	log.Printf("TestDBPositiveCase1 Testing")

	ctxBg := context.Background()

	// variable dipake nanti
	_, err := NewConn(ctxBg, os.Getenv("PG_CONNSTRING"))

	if err != nil {
		// log.Fatalf("error init DB ====>: %s\n", err.Error())
		t.Errorf("errPing must be nil, and should NOT give error, because credential is valid")
	}

}
func TestDBNegativeCase1(t *testing.T) {

	log.Printf("TestDBNegativeCase1 Testing")

	ctxBg := context.Background()

	// variable dipake nanti
	_, err := NewConn(ctxBg, os.Getenv("PG_CONNSTRINGX"))
	// log.Printf("log err: %s", err.Error())
	if err == nil {
		// log.Fatalf("error init DB ====>: %s\n", err.Error())
		t.Errorf("errPing must NOT nil, and should give error, because credential is NOT valid")
	}

}
