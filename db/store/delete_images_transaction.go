package db

import (
	"context"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type DeleteImagesTransactionParams struct {
	AccountID int64                `json:"account_id"`
	ImageIDs  []primitive.ObjectID `json:"image_ids"`
	Amount    int64                `json:"amount"`
}


// Delete Images handles deletion of multiple images from the databases.
func (store *DBStore) DeleteImagesTransaction(ctx context.Context, arg DeleteImagesTransactionParams) (db.Account, error) {
	var uploader db.Account

	err := store.execTransaction(
		ctx,
		func(q *db.Queries) error {
			var err error
			uploader, err = q.UpdateImageCount(ctx, db.UpdateImageCountParams{
				Amount: arg.Amount * -1,
				ID: arg.AccountID,
			})
			return err
		},
		func(op *mongodb.MongoOperations) error {
			err := op.DeleteImages(ctx, arg.ImageIDs)
			return err
		},
	)
	if err != nil {
		return db.Account{}, err
	}

	return uploader, nil
}
