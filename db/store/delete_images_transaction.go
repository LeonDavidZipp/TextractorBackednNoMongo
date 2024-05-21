package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	bucket "github.com/LeonDavidZipp/Textractor/db/s3_bucket"
)


type DeleteImagesTransactionParams struct {
	UserID      int64    `json:"user_id"`
	ImageIDs    []int64  `json:"image_ids"`
	Urls        []string `json:"urls"`
	PreviewUrls []string `json:"preview_urls"`
	Amount      int64    `json:"amount"`
}

// Delete Images handles deletion of multiple images from the databases.
func (store *DBStore) DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.User, error) {
	var uploader db.User

	err := store.execTransaction(
		ctx,
		func(c *bucket.Client) error {
			return c.DeleteImagesFromS3(ctx, bucket.DeleteImagesFromS3Params{
				Urls: arg.Urls,
				PreviewUrls: arg.PreviewUrls,
		})
		},
		func(op *bucket.Client) error {
			return nil
		},
		func(q *db.Queries) error {
			err := q.DeleteImages(ctx, arg.ImageIDs)
			if err != nil {
				return err
			}

			uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: arg.Amount * -1,
				ID: arg.UserID,
				},
			)
			return err
		},
	)
	if err != nil {
		return db.User{}, err
	}

	return uploader, nil
}
