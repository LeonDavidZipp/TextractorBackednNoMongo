package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
)


type UploadImageTransactionParams struct {
	AccountID int64                 `json:"account_id"`
	File      *multipart.FileHeader `json:"file"`
}

type UploadImageTransactionResult struct {
	Image    mongodb.Image `json:"image"`
	Uploader db.Account    `json:"uploader"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *DBStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error) {
	var uploader db.Account
	var image mongodb.Image
	var link string
	var text string

	err := store.execTransaction(
		ctx,
		func(c *bucket.Client) error {
			result, err := c.UploadAndExtractImage(ctx, arg.File)
			if err != nil {
				return err
			}

			link = result.Link
			text = result.Text
			return nil
		},
		func(c *bucket.Client) error {
			return c.DeleteImagesFromS3(ctx, []string{link})
		},
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
				Text: text,
				Link: link,
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
