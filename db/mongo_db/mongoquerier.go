package db

import (
	"context"
)


type MongoQuerier interface {
	InsertImage(ctx context.Context, arg InsertImageParams) (Image, error)
	FindImage(ctx context.Context, id primitive.ObjectID) (Image, error)
	ListImages(ctx context.Context, arg ListParams) ([]Image, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
	DeleteImage(ctx context.Context, id primitive.ObjectID) error
	DeleteImages(ctx context.Context, arg DeleteImagesParams) error
}

var _ MongoQuerier = (*MongoQueries)(nil)
