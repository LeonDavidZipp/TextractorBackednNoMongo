package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
	"mime/multipart"
)


type UploadImageTransactionParams struct {
	UserID int64                 `json:"user_id"`
	Image     *multipart.FileHeader `json:"image"`
}

type UploadImageTransactionResult struct {
	Uploader db.User    `json:"uploader"`
	Image    db.Image   `json:"image"`
}

// Upload Image handles uploading the necessary data and image to the databases.
func (store *DBStore) UploadImageTransaction(ctx context.Context, arg UploadImageTransactionParams) (UploadImageTransactionResult, error) {
	var uploader db.User
	var image db.Image
	var url string
	var previewUrl string
	var text string

	err := store.execTransaction(
		ctx,
		func(c *bucket.Client) error {
			result, err := c.UploadAndExtractImage(ctx, arg.Image)
			if err != nil {
				return err
			}

			url = result.Url
			previewUrl = result.PreviewUrl
			text = result.Text
			return nil
		},
		func(c *bucket.Client) error {
			return c.DeleteImagesFromS3(ctx, []string{url})
		},
		func(q *db.Queries) error {
			var err error
			image, err = q.CreateImage(ctx, db.CreateImageParams{
				UserID: arg.UserID,
				Url: url,
				PreviewUrl: previewUrl,
				Text: text,
			})
			if err != nil {
				return err
			}

			uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: 1,
				ID: arg.UserID,
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
