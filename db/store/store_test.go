package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	"github.com/LeonDavidZipp/Textractor/util"
	"database/sql"
)

var testImage string

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
	)
	account := createRandomAccount(t)

	amount := 10

	results := make(chan UploadImageTransactionResult, amount)
	errs := make(chan error, amount)
	defer close(results)
	defer close(errs)


	for i := 0; i < amount; i++ {
		go func() {
			ctx := context.Background()

			result, err := store.UploadImageTransaction(
				ctx,
				UploadImageTransactionParams{
					AccountID: account.ID,
					Text: "some text",
					Link: "some link",
					Image64: "some image",
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
