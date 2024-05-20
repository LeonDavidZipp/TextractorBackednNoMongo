package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type DeleteImagesTransactionParams struct {
	UserID int64                `json:"user_id"`
	ImageIDs  []primitive.ObjectID `json:"image_ids"`
	URLs     []string             `json:"urls"`
	Amount    int64                `json:"amount"`
}


// Delete Images handles deletion of multiple images from the databases.
func (store *DBStore) DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.User, error) {
	var uploader db.User

	err := store.execTransaction(
		ctx,
		func(c *bucket.Client) error {
			return c.DeleteImagesFromS3(ctx, arg.URLs)
		},
		func(op *bucket.Client) error {
			return nil
		},
		func(q *db.Queries) error {
			var err error
			uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: arg.Amount * -1,
				ID: arg.UserID,
				},
			)
			return err
		},
		func(op *mongodb.MongoOperations) error {
			return op.DeleteImagesFromMongo(ctx, arg.ImageIDs)
		},
	)
	if err != nil {
		return db.User{}, err
	}

	return uploader, nil
}
