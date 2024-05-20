package db

import (
	"context"
	"database/sql"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
)


type Store interface {
	db.Querier
	bucket.S3Client
	UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error)
	DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.User, error)
}

type DBStore struct {
	*db.Queries
	*bucket.Client
	DB *sql.DB
}

func NewStore(DB *sql.DB, s3Client *s3.Client) Store {
	return &DBStore{
		Queries: db.New(DB),
		Client: bucket.NewS3(s3Client),
		DB: DB,
	}
}
