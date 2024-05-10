package db
import (
	"context"
	"database/sql"
	"fmt"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)


type InsertImageParams struct {
	AccountID Account            `bson:"account_id" json:"account_id"`
	Text      string             `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string             `bson:"link" json:"link"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	// base64 encoded image
	Image64   string             `bson:"image_64" json:"image_64"`
}

// InsertImage inserts a new image into the database
// TODO: implement not for Store but for Mongo "Queries"
func (store *Store) InsertImage(
	ctx context.Context,
	arg InsertImageParams
) (InsertImageResult, error) {
	collection := store.ImageDB.Collection("images")
	imageID, err := collection.InsertOne(ctx, arg)
	if err != nil {
		return UploadImageTransactionResult{}, fmt.Errorf("Could not insert image: %w", err)
	}

	result := Image{
		ID:        imageID,
		AccountID: arg.AccountID,
		Text:      arg.Text,
		Link:      arg.Link,
		CreatedAt: arg.CreatedAt,
		Image64:   arg.Image64,
	}

	return result, nil
}

func (store *Store) FindImage(ctx context.Context, id primitive.ObjectID) (Image, error) {
	collection := store.ImageDB.Collection("images")
	filter := bson.M{"_id": id}
	var image Image
	err := collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		return Image{}, fmt.Errorf("Could not find image: %w", err)
	}
	return image, nil
}
