package db

import (
	"context"
	"database/sql"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
)


type Store interface {
	db.Querier
	mongodb.MongoOperator
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error)
	DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.Account, error)
}

type SQLMongoStore struct {
	*db.Queries
	*mongodb.MongoOperations
	UserDB  *sql.DB
	ImageDB *mongo.Database
	ImageBucket *s3.S3
}

func NewStore(userDB *sql.DB, imageDB *mongo.Database) Store {
	return &SQLMongoStore{
		Queries: db.New(userDB),
		MongoOperations: mongodb.NewMongo(imageDB),
		UserDB: userDB,
		ImageDB: imageDB,
	}
}
