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
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var testAccountQueries *db.Queries
var testAccountDB *sql.DB

var testImageOperations *mongodb.MongoOperations
var testImageDB *mongo.Database

var testImageClient *s3.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testAccountDB, err = sql.Open(
		os.Getenv("POSTGRES_DRIVER"),
		os.Getenv("POSTGRES_SOURCE"),
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}
	defer testAccountDB.Close()

	testAccountQueries = db.New(testAccountDB)

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

	testImageOperations = mongodb.NewMongo(testImageDB)

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
