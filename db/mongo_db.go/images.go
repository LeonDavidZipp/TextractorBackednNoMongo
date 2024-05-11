package db
import (
	"context"
	"database/sql"
	"fmt"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)


type InsertImageParams struct {
	AccountID int64              `bson:"account_id" json:"account_id"`
	Text      string             `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string             `bson:"link" json:"link"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Image64   string             `bson:"image_64" json:"image_64"`
}

// InsertImage inserts a new image into the database
// TODO: implement not for Store but for Mongo "Queries"
func (store *Store) InsertImage(
	ctx context.Context,
	arg InsertImageParams
) (Image, error) {
	collection := store.ImageDB.Collection("images")
	inserted, err := collection.insertOne(ctx, arg)

	if err != nil {
		return Image{}, fmt.Errorf("Could not insert image: %w", err)
	}

	image := Image{
		ID:        inserted.imageID,
		AccountID: arg.AccountID,
		Text:      arg.Text,
		Link:      arg.Link,
		CreatedAt: arg.CreatedAt,
		Image64:   arg.Image64,
	}

	return image, nil
}

func (store *Store) FindImage(ctx context.Context, id primitive.ObjectID) (Image, error) {
	collection := store.ImageDB.Collection("images")
	filter := bson.M{"_id": id}

	var image Image
	image, err := collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		return Image{}, fmt.Errorf("Could not find image: %w", err)
	}
	return image, nil
}

func (store *Store) DeleteImage(ctx context.Context, id primitive.ObjectID) error {
	collection := store.ImageDB.Collection("images")
	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Could not delete image: %w", err)
	}

	return nil
}

type UpdateImageParams struct {
	ImageID   primitive.ObjectID `bson:"image_id" json:"image_id"`
	Text      string             `bson:"text" json:"text"`
	Image64   string             `bson:"image_64" json:"image_64"`
}

func (store *Store) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	collection := store.ImageDB.Collection("images")
	filter := bson.M("_id": arg.ImageID)

	image, err := collection.updateOne(ctx, filter, arg)
	if err != nil {
		return fmt.Errorf("Could not update image: %w", err)
	}

	return image, nil
}
