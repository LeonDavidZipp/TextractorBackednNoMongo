package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
	"mime/multipart"
)


type UploadImageTransactionParams struct {
	UserID int64                 `json:"user_id"`
	Image     *multipart.FileHeader `json:"image"`
}

type UploadImageTransactionResult struct {
	Uploader db.User    `json:"uploader"`
	Image    mongodb.Image `json:"image"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *DBStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error) {
	var uploader db.User
	var image mongodb.Image
	var link string
	var text string

	err := store.execTransaction(
		ctx,
		func(c *bucket.Client) error {
			result, err := c.UploadAndExtractImage(ctx, arg.Image)
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
				ID: arg.UserID,
			})
			return err
		},
		func(op *mongodb.MongoOperations) error {
			var err error
			image, err = op.InsertImage(ctx, mongodb.InsertImageParams{
				UserID: arg.UserID,
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
