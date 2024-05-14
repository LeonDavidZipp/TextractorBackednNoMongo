package db

import (
	"context"
	"fmt"
	"time"
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
	fmt.Println("Testing upload image transaction...")
	store := NewStore(
		testAccountDB,
		testImageDB,
	)
	account := createRandomAccount(t)

	results := make(chan UploadImageTransactionResult)
	errs := make(chan error)

	for i := 0; i < 10; i++ {
		go func(i int) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
			defer cancel()

			result, err := store.UploadImageTransaction(
				ctx,
				UploadImageTransactionParams{
					AccountID: account.ID,
					Text: fmt.Sprintf("transaction text %d", i),
					Link: fmt.Sprintf("transaction link %d", i),
					Image64: fmt.Sprintf("transaction image %d", i),
				},
			)

			results <- result
			errs <- err
		}(i)
	}

	for i := 0; i < 10; i++ {
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
		require.Equal(t, account.ImageCount + 1, uploader.ImageCount)
		require.NotZero(t, image.ID)

		_, err = store.GetAccount(ctx, account.ID)
		require.NoError(t, err)
		_, err = store.FindImage(ctx, image.ID)
		require.NoError(t, err)
	}
}
