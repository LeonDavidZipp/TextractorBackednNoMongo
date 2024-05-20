package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
	"database/sql"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
	mongo
	"github.com/LeonDavidZipp/Textractor/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func createRandomUser(t *testing.T) db.User {
	arg := db.CreateUserParams{
		Owner : util.RandomName(),
		Email : util.RandomEmail(),
		GoogleID : sql.NullString{},
		FacebookID : sql.NullString{},
	}

	ctx := context.Background()
	user, err := testUserQueries.CreateUser(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Owner, user.Owner)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.GoogleID, user.GoogleID)
	require.Equal(t, arg.FacebookID, user.FacebookID)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestUploadImageTransaction(t *testing.T) {
	store := NewStore(
		testUserDB,
		testImageDB,
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
		_, err = store.FindImage(ctx, image.ID)
		require.NoError(t, err)
	}

	updatedUser, err := store.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ImageCount + int64(amount), updatedUser.ImageCount)
}

func TestDeleteImagesTransaction(t *testing.T) {
	store := NewStore(
		testUserDB,
		testImageDB,
		testImageClient,
	)
	ctx := context.Background()
	
	user, err := store.CreateUser(
		ctx,
		db.CreateUserParams{
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
				UserID: user.ID,
				Text: "some text",
				Link: "some link",
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
