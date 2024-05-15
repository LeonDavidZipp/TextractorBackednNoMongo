package db

import (
	"context"
	"database/sql"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/mongo"
)


type Store interface {
	db.Querier
	mongodb.MongoOperator
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error)
}

type SQLMongoStore struct {
	*db.Queries
	*mongodb.MongoOperations
	UserDB  *sql.DB
	ImageDB *mongo.Database
}

func NewStore(userDB *sql.DB, imageDB *mongo.Database) Store {
	return &SQLMongoStore{
		Queries: db.New(userDB),
		MongoOperations: mongodb.NewMongo(imageDB),
		UserDB: userDB,
		ImageDB: imageDB,
	}
}
