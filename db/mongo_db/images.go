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
func (q *MongoQueries) InsertImage(
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

func (q *MongoQueries) FindImage(ctx context.Context, id primitive.ObjectID) (Image, error) {
	collection := store.ImageDB.Collection("images")
	filter := bson.M{"_id": id}

	var image Image
	image, err := collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		return Image{}, fmt.Errorf("Could not find image: %w", err)
	}
	return image, nil
}

func (q *MongoQueries) DeleteImage(ctx context.Context, id primitive.ObjectID) error {
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

func (q *MongoQueries) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	collection := store.ImageDB.Collection("images")
	filter := bson.M{"_id": arg.ImageID}
	update := bson.M{
		"$set": bson.M{
			"text": arg.Text,
			"image64": arg.Image64
		}
	}
	options := {returnDocument: "after"}

	result, err := collection.FindOneAndUpdate(ctx,
		filter,
		update,
		options)
	if err != nil {
		return fmt.Errorf("Could not update image: %w", err)
	}

	return image, nil
}

type ListImagesParams struct {
	AccountID int64 `bson:"account_id" json:"account_id"`
	Amount    int32 `bson:"amount" json:"amount"`
	Offset    int32 `bson:"offset" json:"offset"`
}

func (q *MongoQueries) ListImages(ctx context.Context, arg ListParams) ([]Image, error) {

}

type DeleteImagesParams struct {
	AccountIDs []int64 `bson:"account_ids" json:"account_ids"`
	Amount     int32 `bson:"amount" json:"amount"`
	Offset     int32 `bson:"offset" json:"offset"`
}

func (q *MongoQueries) DeleteImages(ctx context.Context, arg ListParams) ([]Image, error) {

}
