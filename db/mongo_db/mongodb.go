package db

import (
	// "context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// type MongoDBTX interface {
// 	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
// 	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
// 	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
// 	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
// 	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
// }

// type MongoDBTX interface {
// 	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
// }

type MongoDBSession interface {
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

type MongoDBCollection interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
	FindOneAndUpdate(context.Context, interface{}, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
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

// func NewMongo(db MongoDBTX) *MongoOperations {
// 	return &MongoOperations{db: db}
// }

// type MongoOperations struct {
// 	db MongoDBTX
// }

// func (q *MongoOperations) WithTx(tx *mongo.Client) *MongoOperations {
// 	return &MongoOperations{
// 		db: tx,
// 	}
// }
