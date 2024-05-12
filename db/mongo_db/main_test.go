package db

import (
	_ "github.com/lib/pq"
	"context"
	"os"
	"log"
	"testing"
	"github.com/LeonDavidZipp/Textractor/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var testImageQueries *MongoOperations
var testImageDB *mongo.Database

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	ctx := context.Background()

	optionsClient := options.Client().ApplyURI(config.MongoSource)
	mongoClient, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("Cannot connect to Image DB:", err)
	}
	defer mongoClient.Disconnect(ctx)

	testImageDB = mongoClient.Database(config.MongoDBName)
	err = testImageDB.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatal("Image DB not reachable:", err)
	}

	testImageQueries = NewMongo(testImageDB)

	os.Exit(m.Run())
}
