package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	"database/sql"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"github.com/LeonDavidZipp/Textractor/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func createRandomAccount(t *testing.T) db.Account {
	arg := db.CreateAccountParams{
		Owner : util.RandomName(),
		Email : util.RandomEmail(),
		GoogleID : sql.NullString{},
		FacebookID : sql.NullString{},
	}

	ctx := context.Background()
	account, err := testAccountQueries.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.GoogleID, account.GoogleID)
	require.Equal(t, arg.FacebookID, account.FacebookID)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestUploadImageTransaction(t *testing.T) {
	store := NewStore(
		testAccountDB,
		testImageDB,
		testImageClient,
	)

	account := createRandomAccount(t)

	image, err := util.ImageAsFileHeader("/app/test_files/sample.jpeg")
	require.NoError(t, err)

	amount := 2

	results := make(chan UploadImageTransactionResult, amount)
	errs := make(chan error, amount)

	for i := 0; i < amount; i++ {
		go func() {
			ctx := context.Background()

			result, err := store.UploadImageTransaction(
				ctx,
				UploadImageTransactionParams{
					AccountID: account.ID,
					Image: image,
				},
			)

			results <- result
			errs <- err
		}()
	}

	for i := 0; i < amount; i++ {
		ctx := context.Background()
		err := <- errs
		require.NoError(t, err)

		result := <-results
		uploader := result.Uploader
		image := result.Image

		// check uploader account and account are the same
		require.NotEmpty(t, uploader)
		require.NotEmpty(t, image)

		require.Equal(t, account.ID, uploader.ID)
		require.Equal(t, uploader.ID, image.AccountID)
		require.NotZero(t, image.ID)
		
		_, err = store.GetAccount(ctx, account.ID)
		require.NoError(t, err)
		_, err = store.FindImage(ctx, image.ID)
		require.NoError(t, err)
	}

	updatedAccount, err := store.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ImageCount + int64(amount), updatedAccount.ImageCount)
}

func TestDeleteImagesTransaction(t *testing.T) {
	store := NewStore(
		testAccountDB,
		testImageDB,
		testImageClient,
	)
	ctx := context.Background()
	
	account, err := store.CreateAccount(
		ctx,
		db.CreateAccountParams{
			Owner: util.RandomName(),
			Email: util.RandomEmail(),
		},
	)
	require.NoError(t, err)

	amount := 2
	imageIDs := make([]primitive.ObjectID, amount)
	for i := 0; i < amount; i++ {
		image, err := store.InsertImage(
			ctx,
			mongodb.InsertImageParams{
				AccountID: account.ID,
				Text: "some text",
				Link: "some link",
			},
		)
		require.NoError(t, err)
		imageIDs[i] = image.ID
	}

	toDelete := imageIDs[:amount/2]
	updatedAccount, err := store.DeleteImagesTransaction(
		ctx,
		DeleteImagesTransactionParams{
			AccountID: account.ID,
			ImageIDs: toDelete,
			Amount: int64(len(toDelete)),
		},
	)
	require.NoError(t, err)
	require.Equal(t, account.ImageCount - int64(len(toDelete)), updatedAccount.ImageCount)

	for _, id := range toDelete {
		image, err := store.FindImage(ctx, id)
		require.Error(t, err)
		require.Empty(t, image)
	}

	remaining := imageIDs[amount/2:]
	for _, id := range remaining {
		image, err := store.FindImage(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, image)
	}
}
