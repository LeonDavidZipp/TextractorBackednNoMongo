package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"log"
	"testing"
	"go.mongodb.org/mongo-driver/mongo/options"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongo
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var testQueries *db.Queries
var testDB *sql.DB

var testImageClient *s3.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testDB, err = sql.Open(
		os.Getenv("POSTGRES_DRIVER"),
		os.Getenv("POSTGRES_SOURCE"),
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}
	defer testDB.Close()

	testQueries = db.New(testDB)

	config, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal("Cannot load AWS config:", err)
	}

	testImageClient = s3.NewFromConfig(config)

	os.Exit(m.Run())
}
