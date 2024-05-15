package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type InsertImageParams struct {
	AccountID int64  `bson:"account_id" json:"account_id"`
	Text      string `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string `bson:"link" json:"link"`
	Image64   string `bson:"image_64" json:"image_64"`
}

// InsertImage inserts a new image into the database
// TODO: implement not for Store but for Mongo "Queries"
func (op *MongoOperations) InsertImage(ctx context.Context, arg InsertImageParams) (Image, error) {
	inserted, err := op.db.Collection("imagedb").InsertOne(ctx, arg)

	if err != nil {
		return Image{}, fmt.Errorf("Could not insert image: %w", err)
	}

	image := Image{
		ID:        inserted.InsertedID.(primitive.ObjectID),
		AccountID: arg.AccountID,
		Text:      arg.Text,
		Link:      arg.Link,
		Image64:   arg.Image64,
	}

	return image, nil
}

func (op *MongoOperations) FindImage(ctx context.Context, id primitive.ObjectID) (Image, error) {
	filter := bson.M{"_id": id}

	var image Image
	err := op.db.Collection("imagedb").FindOne(ctx, filter).Decode(&image)
	if err != nil {
		return Image{}, fmt.Errorf("Could not find image: %w", err)
	}
	return image, nil
}

type ListImagesParams struct {
	AccountID int64 `bson:"account_id" json:"account_id"`
	Limit     int64 `bson:"amount" json:"amount"`
	Offset    int64 `bson:"offset" json:"offset"`
}

func (op *MongoOperations) ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error) {
	filter := bson.M{"account_id": arg.AccountID}

	findOptions := options.Find()
	findOptions.SetLimit(arg.Limit)
	findOptions.SetSkip(arg.Offset)

	cursor, err := op.db.Collection("imagedb").Find(ctx, filter, findOptions)
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
}

func (op *MongoOperations) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	filter := bson.M{"_id": arg.ImageID}
	update := bson.M{
		"$set": bson.M{
			"text": arg.Text,
		},
	}
	findOptions := options.FindOneAndUpdate().SetReturnDocument(options.After)

	result := op.db.Collection("imagedb").FindOneAndUpdate(ctx,
		filter,
		update,
		findOptions,
	)

	var image Image
	err := result.Decode(&image)
	if err != nil {
		return Image{}, err
	}
	return image, nil
}

func (op *MongoOperations) DeleteImage(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := op.db.Collection("imagedb").DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Could not delete image: %w", err)
	}

	return nil
}

func (op *MongoOperations) DeleteImages(ctx context.Context, ids []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err := op.db.Collection("imagedb").DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("Could not delete images: %w", err)
	}

	return nil
}
