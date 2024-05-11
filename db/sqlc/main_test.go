package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"log"
	"testing"
	"github.com/LeonDavidZipp/Textractor/util"
)

var testUserQueries *Queries
var testUserDB *sql.DB

var testImageQueries *MongoQueries
var testtestImageDB *mongo.Database

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	testUserDB, err = sql.Open(
		config.DBDriver,
		config.DBSource,
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}
	testUserQueries = New(testUserDB)

	optionsClient := options.Client().ApplyURI(config.MongoSource)
	mongoClient, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("Cannot connect to Image DB:", err)
	}
	defer mongoClient.Disconnect(ctx)

	testImageDB := mongoClient.Database(config.MongoDBName)
	err := testImageDB.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatal("Image DB not reachable:", err)
	}

	testImageQueries = NewMongo(testImageDB)

	os.Exit(m.Run())
}
