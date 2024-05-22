package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	"github.com/LeonDavidZipp/Textractor/util"
)


func createRandomUser(t *testing.T) db.User {
	name := util.RandomString(8)

	ctx := context.Background()
	user, err := testQueries.CreateUser(ctx, name)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, name, user.Name)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestUploadImageTransaction(t *testing.T) {
	store := NewStore(
		testDB,
		testImageClient,
	)

	user := createRandomUser(t)

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
					UserID: user.ID,
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

		// check uploader user and user are the same
		require.NotEmpty(t, uploader)
		require.NotEmpty(t, image)

		require.Equal(t, user.ID, uploader.ID)
		require.Equal(t, uploader.ID, image.UserID)
		require.NotZero(t, image.ID)
		
		_, err = store.GetUser(ctx, user.ID)
		require.NoError(t, err)
		_, err = store.GetImageFromSQL(ctx, image.ID)
		require.NoError(t, err)
	}

	updatedUser, err := store.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ImageCount + int64(amount), updatedUser.ImageCount)
}

func TestDeleteImagesTransaction(t *testing.T) {
	store := NewStore(
		testDB,
		testImageClient,
	)
	ctx := context.Background()
	
	user, err := store.CreateUser(
		ctx,
		util.RandomString(8),
	)
	require.NoError(t, err)

	amount := 2
	imageIDs := make([]int64, amount)
	for i := 0; i < amount; i++ {
		image, err := store.CreateImage(
			ctx,
			 db.CreateImageParams{
				UserID: user.ID,
				Url: "some url",
				PreviewUrl: "some preview url",
				Text: "some text",
			},
		)
		require.NoError(t, err)
		imageIDs[i] = image.ID
	}

	toDelete := imageIDs[:amount/2]
	updatedUser, err := store.DeleteImagesTransaction(
		ctx,
		DeleteImagesTransactionParams{
			UserID: user.ID,
			ImageIDs: toDelete,
			Amount: int64(len(toDelete)),
		},
	)
	require.NoError(t, err)
	require.Equal(t, user.ImageCount - int64(len(toDelete)), updatedUser.ImageCount)

	for _, id := range toDelete {
		image, err := store.GetImageFromSQL(ctx, id)
		require.Error(t, err)
		require.Empty(t, image)
	}

	remaining := imageIDs[amount/2:]
	for _, id := range remaining {
		image, err := store.GetImageFromSQL(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, image)
	}
}
