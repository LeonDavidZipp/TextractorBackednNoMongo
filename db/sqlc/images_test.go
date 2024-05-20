package db

import (
	"context"
	"time"
	"testing"
	"database/sql"
)


func createRandomImage(t *testing.T, userID int64) Image {
	url := util.RandomString(8)
	text := util.RandomText()

	ctx := context.Background()
	image, err := testImageQueries.CreateImage(ctx, CreateImageParams{
		UserID: userID,
		Url: url,
		Text: text,
	})
	require.NoError(t, err)
	require.NotEmpty(t, image)
	require.Equal(t, user.ID, image.UserID)
	require.Equal(t, url, image.Url)
	require.Equal(t, text, image.Text)

	require.NotZero(t, image.ID)
	require.NotZero(t, image.CreatedAt)

	return image
}

func TestCreateImage(t *testing.T) {
	createRandomImage(t)
}

func TestGetImage(t *testing.T) {
	user := createRandomUser(t)
	image1 := createRandomImage(t, user.ID)

	ctx := context.Background()
	image2, err := testImageQueries.GetImage(ctx, image1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.ID, image2.ID)
	require.Equal(t, image1.UserID, image2.UserID)
	require.Equal(t, image1.Url, image2.Url)
	require.Equal(t, image1.Text, image2.Text)

	require.WithinDuration(t, image1.CreatedAt, image2.CreatedAt, 10 * time.Second)
}

func TestListImages(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomImage(t, user.ID)
	}

	ctx := context.Background()
	images, err := testImageQueries.ListImages(ctx)

	require.NoError(t, err)
	require.Len(t, images, 5)

	for _, image := range images {
		require.NotEmpty(t, image)
		require.NotZero(t, image.ID)
		require.Equal(t, user.ID, image.UserID)
		require.NotEmpty(t, image.Url)
		require.NotEmpty(t, image.Text)
		require.NotZero(t, image.CreatedAt)
	}
}

func TestDeleteImages(t *testing.T) {
	user := createRandomUser(t)

	var image Image
	var imageIDs []int32
	for i := 0; i < 10; i++ {
		image = createRandomImage(t, user.ID)
		imageIDs = append(imageIDs, int32(image.ID))
	}

	ctx := context.Background()

	err := testImageQueries.DeleteImages(ctx, imageIDs)
	require.NoError(t, err)

	for _, id := range imageIDs {
		image, err = testImageQueries.GetImage(ctx, id)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows)
		require.Empty(t, image)
	}
}
