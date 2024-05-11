package db

import (
	"context"
	"fmt"
	"os"
	"log"
	"encoding/base64"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func encodeImageToBase64(filepath string) string {
	imageData, err := ioutil.ReadFile(filepath)
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	return base64Image
}

exampleImage1 := encodeImageToBase64("/Users/leon/Desktop/Textractor/db/mongo_db/sample.jpeg")
exampleImage2 := encodeImageToBase64("/Users/leon/Desktop/Textractor/db/mongo_db/text.png")

func insertImage(t *testing.T, image64 string, accountID int64) Image {
	arg := InsertImageParams{
		AccountID: accountID,
		Text: RandomString(100),
		Link: RandomLink(),
		Image64: exampleImage1,
	}

	ctx := context.Background()
	image, err := testImageQueries.InsertImage(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, image)

	require.Equal(t, arg.AccountID, image.AccountID)
	require.Equal(t, arg.Text, image.Text)
	require.Equal(t, arg.Link, image.Link)
	require.Equal(t, arg.Image64, image.Image64)

	return image
}

func TestInsertImage(t *testing.T) {
	insertImage(t, exampleImage1, 1)
}

func TestFindImage(t *testing.T) {
	image1 := insertImage(t, exampleImage1, 1)
	image2, err := FindImage(ctx, image1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.AccountID, image2.AccountID)
	require.Equal(t, image1.Text, image2.Text)
	require.Equal(t, image1.Link, image2.Link)
	require.Equal(t, image1.Image64, image2.Image64)
}

func TestListImages(t *testing.T) {
	for i := 0; i < 10; i++ {
		insertImage(t, exampleImage1, 1)
	}

	arg := ListImagesParams{
		AccountID: 1,
		Amount: 5,
		Offset: 5,
	}

	ctx := context.Background()
	images, err := testImageQueries.ListImages(ctx, arg)

	require.NoError(t, err)
	require.Len(t, images, 5)

	for _, image := range images {
		require,NotEmpty(t, image)
	}
}

func TestUpdateImage(t *testing.T) {
	image1 := insertImage(t, exampleImage1, 1)
	arg := UpdateImageParams{
		ID: image1.ID,
		Text: RandomText(),
	}

	ctx := context.Background()
	image2, err := testImageQueries.UpdateImage(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, image2)

	require.Equal(t, image1.ID, image2.ID)
	require.Equal(t, arg.Text, image2.Text)

	require.Equal(t, image1.AccountID, image2.AccountID)
	require.Equal(t, image1.Link, image2.Link)
	require.Equal(t, image1.Image64, image2.Image64)
}

func TestDeleteImage(t *testing.T) {
	image1 := insertImage(t, exampleImage1, 1)

	ctx := context.Background()
	err := testImageQueries.DeleteImage(ctx, image1.ID)
	require.NoError(t, err)

	image2, err := testImageQueries.FindImage(ctx, image1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, image2)
}

func TestDeleteImages(t *testing.T) {
	for i := 0; i < 10; i++ {
		insertImage(t, exampleImage1, 1)
	}

	ctx := context.Background()
	err := testImageQueries.DeleteImages(ctx, 1)
	require.NoError(t, err)

	images, err := testImageQueries.ListImages(ctx, ListImagesParams{AccountID: 1})
	require.Error(t, err)
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
	require.Empty(t, images)
}