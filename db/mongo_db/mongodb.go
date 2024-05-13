package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// type MongoSession interface {
// 	EndSession(context.Context) error
// 	StartTransaction(...*options.TransactionOptions) error
// 	AbortTransaction(context.Context) error
// 	CommitTransaction(context.Context) error
// 	WithTransaction(context.Context, func(mongo.SessionContext) error) error
// }

type MongoDBCollection interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	FindOneAndUpdate(context.Context, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

func NewMongo(db MongoDatabase) *MongoOperations {
	return &MongoOperations{
		db: db,
	}
}

type MongoOperations struct {
	db MongoDatabase
}
