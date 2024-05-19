package db

import (
	"github.com/LeonDavidZipp/Textractor/util"
	"context"
	"testing"
	"github.com/stretchr/testify/require"

	"fmt"
)


// /Users/lzipp/Desktop/Textractor/Backend/db/s3_bucket/images_test.go
// /Users/lzipp/Desktop/Textractor/Backend/test_files/sample.jpeg

func uploadImage(t *testing.T, imagePath string) string {
	ctx := context.Background()

	image, err := util.ImageAsFileHeader(imagePath)
	require.NoError(t, err)
	fmt.Println("\n\n\nImage: ", image)

	result, err := testImageClient.UploadAndExtractImage(ctx, image)
	require.NoError(t, err)

	link := result.Link
	text := result.Text

	require.NotEmpty(t, link)
	require.NotEmpty(t, text)

	return link
}

func TestUploadImage(t *testing.T) {
	uploadImage(t, "/app/test_files/sample.jpeg")
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()

	link := uploadImage(t, "/app/test_files/sample.jpeg")

	imageBytes, err := testImageClient.GetImage(ctx, link)
	require.NoError(t, err)
	require.NotEmpty(t, imageBytes)
}

func TestDeleteImages(t *testing.T) {
	ctx := context.Background()

	links := make([]string, 10)
	for i := 0; i < 10; i++ {
		links[i] = uploadImage(t, "/app/test_files/sample.jpeg")
	}

	err := testImageClient.DeleteImagesFromS3(ctx, links)
	require.NoError(t, err)

	for _, link := range links {
		imageBytes, err := testImageClient.GetImage(ctx, link)
		require.Error(t, err)
		require.Empty(t, imageBytes)
	}
}
