package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/LeonDavidZipp/Textractor/util"
	"testing"
	"github.com/stretchr/testify/require"
)


var exampleImage1 string
var exampleImage2 string

func insertImage(t *testing.T, userID int64) Image {
	arg := InsertImageParams{
		UserID: userID,
		Text: util.RandomString(100),
		Link: util.RandomLink(),
	}

	ctx := context.Background()
	image, err := testImageOperations.InsertImage(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, image)

	require.Equal(t, arg.UserID, image.UserID)
	require.Equal(t, arg.Text, image.Text)
	require.Equal(t, arg.Link, image.Link)

	return image
}

func TestInsertImage(t *testing.T) {
	insertImage(t, 1)
}

func TestFindImage(t *testing.T) {
	ctx := context.Background()

	image1 := insertImage(t, 1)
	image2, err := testImageOperations.FindImage(ctx, image1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.UserID, image2.UserID)
	require.Equal(t, image1.Text, image2.Text)
	require.Equal(t, image1.Link, image2.Link)
}

func TestListImages(t *testing.T) {
	for i := 0; i < 10; i++ {
		insertImage(t, 1)
	}

	arg := ListImagesParams{
		UserID: 1,
		Limit: 5,
		Offset: 5,
	}

	ctx := context.Background()
	images, err := testImageOperations.ListImages(ctx, arg)

	require.NoError(t, err)
	require.Len(t, images, 5)

	for _, image := range images {
		require.NotEmpty(t, image)
	}
}

func TestUpdateImage(t *testing.T) {
	image1 := insertImage(t, 1)
	arg := UpdateImageParams{
		ImageID: image1.ID,
		Text: util.RandomText(),
	}

	ctx := context.Background()
	image2, err := testImageOperations.UpdateImage(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.ID, image2.ID)
	require.Equal(t, arg.Text, image2.Text)

	require.Equal(t, image1.UserID, image2.UserID)
	require.Equal(t, image1.Link, image2.Link)
}

func TestDeleteImage(t *testing.T) {
	image1 := insertImage(t, 1)

	ctx := context.Background()
	err := testImageOperations.DeleteImage(ctx, image1.ID)
	require.NoError(t, err)

	image2, err := testImageOperations.FindImage(ctx, image1.ID)
	require.Error(t, err)
	require.Empty(t, image2)
}

func TestDeleteImages(t *testing.T) {
	image1 := insertImage(t, 1)
	image2 := insertImage(t, 1)

	ctx := context.Background()
	err := testImageOperations.DeleteImagesFromMongo(ctx, []primitive.ObjectID{image1.ID, image2.ID})
	require.NoError(t, err)

	image1Del, err := testImageOperations.FindImage(ctx, image1.ID)
	require.Error(t, err)
	require.Empty(t, image1Del)
	image2Del, err := testImageOperations.FindImage(ctx, image2.ID)
	require.Error(t, err)
	require.Empty(t, image2Del)
}
