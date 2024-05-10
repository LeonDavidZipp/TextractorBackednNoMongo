package main

import (
	"database/sql"
	"log"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	sqlConnection, err := sql.Open(
		config.DBDriver,
		config.DBSource,
	)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	optionsClient := options.Client().ApplyURI(config.Mongo_Source)
	mongoClient, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("Cannot connect to mongo:", err)
	}
	defer mongoClient.Disconnect(ctx)



	store := db.NewStore(sqlConnection, mongoClient)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}