package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
)


type UploadImageTransactionParams struct {
	AccountID int64  `json:"account_id"`
	Text      string `json:"text"`
	Link      string `json:"link"`
	Image64   string `json:"image_64"`
}

type UploadImageTransactionResult struct {
	Image    mongodb.Image `json:"image"`
	Uploader db.Account    `json:"uploader"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *SQLMongoStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error) {
	var uploader db.Account
	var image mongodb.Image

	err := store.execTransaction(
		ctx,
		func(q *db.Queries) error {
			var err error
			uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: 1,
				ID: arg.AccountID,
			})
			return err
		},
		func(op *mongodb.MongoOperations) error {
			var err error
			image, err = op.InsertImage(ctx, mongodb.InsertImageParams{
				AccountID: arg.AccountID,
				Text: arg.Text,
				Link: arg.Link,
				Image64: arg.Image64,
			})
			return err
		},
	)
	if err != nil {
		return UploadImageTransactionResult{}, err
	}

	result := UploadImageTransactionResult{
		Image: image,
		Uploader: uploader,
	}

	return result, nil
}
