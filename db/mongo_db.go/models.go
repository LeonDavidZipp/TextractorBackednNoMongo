package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// this struct is used to store the account information
// only for learning, not most efficient way to store
type Image struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	AccountID Account `bson:"account_id" json:"account_id"`
	Text      string `bson:"text" json:"text"`
	Link      string `bson:"link" json:"link"` // link to the image in s3 storage
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Image64   string `bson:"image_64" json:"image_64"`  // base64 encoded image
}
