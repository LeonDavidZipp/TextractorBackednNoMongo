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
	bucket.S3Client
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error)
	DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.User, error)
}

type DBStore struct {
	*db.Queries
	*mongodb.MongoOperations
	*bucket.Client
	UserDB  *sql.DB
	ImageDB *mongo.Database
}

func NewStore(userDB *sql.DB, imageDB *mongo.Database, s3Client *s3.Client) Store {
	return &DBStore{
		Queries: db.New(userDB),
		MongoOperations: mongodb.NewMongo(imageDB),
		Client: bucket.NewS3(s3Client),
		UserDB: userDB,
		ImageDB: imageDB,
	}
}
