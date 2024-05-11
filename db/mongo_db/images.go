package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type InsertImageParams struct {
	AccountID int64              `bson:"account_id" json:"account_id"`
	Text      string             `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string             `bson:"link" json:"link"`
	Image64   string             `bson:"image_64" json:"image_64"`
}

// InsertImage inserts a new image into the database
// TODO: implement not for Store but for Mongo "Queries"
func (q *MongoQueries) InsertImage(ctx context.Context, arg InsertImageParams) (Image, error) {
	collection := q.db.Collection("images")
	inserted, err := collection.insertOne(ctx, arg)

	if err != nil {
		return Image{}, fmt.Errorf("Could not insert image: %w", err)
	}

	image := Image{
		ID:        inserted.imageID,
		AccountID: arg.AccountID,
		Text:      arg.Text,
		Link:      arg.Link,
		Image64:   arg.Image64,
	}

	return image, nil
}

func (q *MongoQueries) FindImage(ctx context.Context, id primitive.ObjectID) (Image, error) {
	collection := q.db.Collection("images")
	filter := bson.M{"_id": id}

	var image Image
	image, err := collection.FindOne(ctx, filter).Decode(&image)
	if err != nil {
		return Image{}, fmt.Errorf("Could not find image: %w", err)
	}
	return image, nil
}

type ListImagesParams struct {
	AccountID int64 `bson:"account_id" json:"account_id"`
	Amount    int32 `bson:"amount" json:"amount"`
	Offset    int32 `bson:"offset" json:"offset"`
}

func (q *MongoQueries) ListImages(ctx context.Context, arg ListParams) ([]Image, error) {
	collection := q.db.Collection("images")
	filter := bson.M{"account_id": arg.AccountID}

	findOptions := options.Find()
	findOptions.SetLimit(arg.Limit)
	findOptions.SetSkip(arg.Offset)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("Could not list images: %w", err)
	}
	defer cursor.Close(ctx)

	var images []Image
	if err = cursor.All(ctx, &images); err != nil {
		return nil, fmt.Errorf("Could not decode images: %w", err)
	}

	return images, nil
}

type UpdateImageParams struct {
	ImageID   primitive.ObjectID `bson:"image_id" json:"image_id"`
	Text      string             `bson:"text" json:"text"`
	Image64   string             `bson:"image_64" json:"image_64"`
}

func (q *MongoQueries) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	collection := q.db.Collection("images")
	filter := bson.M{"_id": arg.ImageID}
	update := bson.M{
		"$set": bson.M{
			"text": arg.Text,
			"image64": arg.Image64
		}
	}
	findOptions := {returnDocument: "after"}

	result, err := collection.FindOneAndUpdate(ctx,
		filter,
		update,
		findOptions,
	)
	if err != nil {
		return fmt.Errorf("Could not update image: %w", err)
	}

	var image Image
	err := result.Decode(&image)
	if err != nil {
		return Image{}, err
	}
	return image, nil
}

func (q *MongoQueries) DeleteImage(ctx context.Context, id primitive.ObjectID) error {
	collection := q.db.Collection("images")
	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Could not delete image: %w", err)
	}

	return nil
}

type DeleteImagesParams struct {
	IDs    []int64 `bson:"ids" json:"ids"`
	Amount int32   `bson:"amount" json:"amount"`
	Offset int32   `bson:"offset" json:"offset"`
}

func (q *MongoQueries) DeleteImages(ctx context.Context, arg DeleteImagesParams) error {
	collection := q.db.Collection("images")
	filter := bson.M{"_id": bson.M{"$in": arg.IDs}}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("Could not delete images: %w", err)
	}

	return nil
}
