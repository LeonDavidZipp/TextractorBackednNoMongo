package api

import (
	_ "github.com/lib/pq"
	"context"
	"os"
	"time"
	"log"
	"testing"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	st "github.com/LeonDavidZipp/Textractor/db/store"
)


func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	userDB, err := sql.Open(
		os.Getenv("POSTGRES_DRIVER"),
		os.Getenv("POSTGRES_SOURCE"),
	)
	if err != nil {
		log.Fatal("Cannot connect to User DB:", err)
	}

	optionsClient := options.Client().ApplyURI(os.Getenv("MONGO_SOURCE"))
	mongoClient, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("Cannot connect to Image DB:", err)
	}
	defer mongoClient.Disconnect(ctx)

	imageDB := mongoClient.Database(os.Getenv("MONGO_DB_NAME"))
	err = imageDB.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatal("Image DB not reachable:", err)
	}

	store := st.NewStore(userDB, imageDB)
	server := NewServer(store)

	err = server.Start(os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
