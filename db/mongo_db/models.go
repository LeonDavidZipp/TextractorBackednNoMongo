package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// this struct is used to store the user information
// only for learning, not most efficient way to store
type Image struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID int64              `bson:"user_id" json:"user_id"`
	Text      string             `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string             `bson:"link" json:"link"`
}
