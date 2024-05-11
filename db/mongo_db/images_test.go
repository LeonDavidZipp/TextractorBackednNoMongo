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

func insertImage(t *testing.T) Image {
	arg := InsertImageParams{
		AccountID: 4269,
		Text: "Example text",
		Link: "www.example.com",
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

}
