package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// this struct is used to store the account information
// only for learning, not most efficient way to store
type Image struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	AccountID int64              `bson:"account_id" json:"account_id"`
	Text      string             `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string             `bson:"link" json:"link"`
}
