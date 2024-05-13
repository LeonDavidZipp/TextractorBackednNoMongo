package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"log"
	"testing"
)

var testAccountQueries *Queries
var testAccountDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testAccountDB, err = sql.Open(
		os.Getenv("POSTGRES_DRIVER"),
		os.Getenv("POSTGRES_SOURCE"),
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}

	testAccountQueries = New(testAccountDB)

	os.Exit(m.Run())
}
