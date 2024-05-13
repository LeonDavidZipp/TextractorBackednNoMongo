package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MongoDBSession interface {
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

type MongoDBCollection interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	FindOneAndUpdate(context.Context, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

func NewMongo(se MongoDBSession, db MongoDBCollection) *MongoOperations {
	return &MongoOperations{
		se: se,
		db: db,
	}
}

type MongoOperations struct {
	se MongoDBSession
	db MongoDBCollection
}
