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
	config, err := util.LoadConfig("../..")
	testAccountDB, err = sql.Open(
		config.DBDriver,
		config.DBSource,
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}

	testAccountQueries = New(testAccountDB)

	os.Exit(m.Run())
}
