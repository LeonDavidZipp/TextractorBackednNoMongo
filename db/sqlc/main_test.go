package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"log"
	"testing"
	"github.com/LeonDavidZipp/Textractor/util"
)

var testAccountQueries *Queries
var testAccountDB *sql.DB

func TestMain(m *testing.M) {
	testAccountDB, err = sql.Open(
		os.getenv("POSTGRES_DRIVER"),
		os.getenv("POSTGRES_SOURCE"),
	)

	testAccountQueries = New(testAccountDB)

	os.Exit(m.Run())
}
