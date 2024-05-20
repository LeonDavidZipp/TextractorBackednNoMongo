package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"log"
	"testing"
)

var testUserQueries *Queries
var testUserDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testUserDB, err = sql.Open(
		os.Getenv("POSTGRES_DRIVER"),
		os.Getenv("POSTGRES_SOURCE"),
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}
	defer testUserDB.Close()

	testUserQueries = New(testUserDB)

	os.Exit(m.Run())
}
