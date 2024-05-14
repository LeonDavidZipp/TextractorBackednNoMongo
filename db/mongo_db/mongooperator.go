package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type MongoOperator interface {
	InsertImage(ctx context.Context, arg InsertImageParams) (Image, error)
	FindImage(ctx context.Context, id primitive.ObjectID) (Image, error)
	ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
	DeleteImage(ctx context.Context, id primitive.ObjectID) error
	DeleteImages(ctx context.Context, ids []primitive.ObjectID) error
}

var _ MongoOperator = (*MongoOperations)(nil)
