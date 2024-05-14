package db

import (
	_ "github.com/lib/pq"
	"context"
	"os"
	"log"
	"testing"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var testImageOperations *MongoOperations
var testImageDB *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()

	optionsClient := options.Client().ApplyURI(os.Getenv("MONGO_SOURCE"))
	mongoClient, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("Cannot connect to Image DB:", err)
	}
	defer mongoClient.Disconnect(ctx)

	testImageDB = mongoClient.Database(os.Getenv("MONGO_DB_NAME"))
	err = testImageDB.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatal("Image DB not reachable:", err)
	}

	testImageOperations = NewMongo(testImageDB)

	os.Exit(m.Run())
}
