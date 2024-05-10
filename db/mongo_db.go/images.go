package db
import (
	"context"
	"database/sql"
	"fmt"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)


type InsertImageParams struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	AccountID Account            `bson:"account_id" json:"account_id"`
	Text      string             `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string             `bson:"link" json:"link"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	// base64 encoded image
	Image64   string             `bson:"image_64" json:"image_64"`
}

type InsertImageResult struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	AccountID Account            `bson:"account_id" json:"account_id"`
	Text      string             `bson:"text" json:"text"`
	Image64   string             `bson:"image_64" json:"image_64"`
}

func (store *Store) InsertImage(
	ctx context.Context,
	arg InsertImageParams
) (InsertImageResult, error) {
	collection := store.ImageDB.Collection("images")
	imageID, err := collection.InsertOne(ctx, arg)
	if err != nil {
		return UploadImageTransactionResult{}, fmt.Errorf("Could not insert image: %w", err)
	}

	result := InsertImageResult{
		ID:        imageID.InsertedID.(primitive.ObjectID),
		AccountID: arg.AccountID,
		Text:      arg.Text,
		Image64:   arg.Image64,
	}

	return result, nil
}

func (store *Store) GetImage(ctx context.Context, id string) (*Image, error) {

}
