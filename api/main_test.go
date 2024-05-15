package api

import (
	_ "github.com/lib/pq"
	"context"
	"os"
	"log"
	"testing"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	db "github.com/LeonDavidZipp/Textractor/db"
	"database/sql"
	sqldb "github.com/LeonDavidZipp/Textractor/db/sql_db"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
)


func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	userDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}

	optionsClient := options.Client().ApplyURI(config.MongoSource)
	mongoClient, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("Cannot connect to Image DB:", err)
	}
	defer mongoClient.Disconnect(ctx)

	imageDB := mongoClient.Database(config.MongoDBName)
	err = imageDB.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatal("Image DB not reachable:", err)
	}

	store := db.NewStore(userDB, imageDB)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
