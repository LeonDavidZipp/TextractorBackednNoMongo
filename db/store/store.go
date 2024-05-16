package db

import (
	"context"
	"database/sql"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
)


type Store interface {
	db.Querier
	mongodb.MongoOperator
	bucket.S3Client
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error)
	DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.Account, error)
}

type DBStore struct {
	*db.Queries
	*mongodb.MongoOperations
	*bucket.Client
	UserDB  *sql.DB
	ImageDB *mongo.Database
	// S3Uploader *manager.Uploader
}

func NewStore(userDB *sql.DB, imageDB *mongo.Database, s3Client *s3.Client) Store {
	return &DBStore{
		Queries: db.New(userDB),
		MongoOperations: mongodb.NewMongo(imageDB),
		UserDB: userDB,
		ImageDB: imageDB,
		Client: s3Uploader,
	}
}
