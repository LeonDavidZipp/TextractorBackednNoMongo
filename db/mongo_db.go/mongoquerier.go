package db

import (
	"context"
)


type MongoQuerier interface {
	InsertImage(ctx context.Context, arg InsertImageParams) (Image, error)
	FindImage(ctx context.Context, id primitive.ObjectID) (Image, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)

}

var _ MongoQuerier = (*MongoQueries)(nil)
