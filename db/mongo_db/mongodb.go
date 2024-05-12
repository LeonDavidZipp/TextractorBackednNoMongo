package db

import (
	// "context"
	"go.mongodb.org/mongo-driver/mongo"
)


// type MongoDBTX interface {
// 	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
// 	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
// 	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
// 	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
// 	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
// }

type MongoDBTX interface {
	Collection(string) *mongo.Collection
}

func NewMongo(db MongoDBTX) *MongoQueries {
	return &MongoQueries{db: db}
}

type MongoQueries struct {
	db MongoDBTX
}

// func (q *MongoQueries) WithTx(tx *mongo.Client) *MongoQueries {
// 	return &MongoQueries{
// 		db: tx,
// 	}
// }
