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

func uploadImage(t *testing.T, imagePath string) (string, string) {
	ctx := context.Background()

	image, err := util.ImageAsFileHeader(imagePath)
	require.NoError(t, err)
	fmt.Println("\n\n\nImage: ", image)

	result, err := testImageClient.UploadAndExtractImage(ctx, image)
	require.NoError(t, err)

	url := result.Url
	previewUrl := result.PreviewUrl
	text := result.Text

	require.NotEmpty(t, url)
	require.NotEmpty(t, previewUrl)
	require.NotEmpty(t, text)

	return url
}

func TestUploadImage(t *testing.T) {
	uploadImage(t, "/app/test_files/sample.jpeg")
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()

	url, _ := uploadImage(t, "/app/test_files/sample.jpeg")

	imageBytes, err := testImageClient.GetImageFromS3(ctx, url)
	require.NoError(t, err)
	require.NotEmpty(t, imageBytes)
}

func TestGetPreview(t *testing.T) {
	ctx := context.Background()

	_, previewUrl := uploadImage(t, "/app/test_files/sample.jpeg")

	previewImageBytes, err := testImageClient.GetPreviewFromS3(ctx, previewUrl)
	require.NoError(t, err)
	require.NotEmpty(t, previewImageBytes)
}

func TestDeleteImages(t *testing.T) {
	ctx := context.Background()

	urls := make([]string, 10)
	previewUrls := make([]string, 10)
	for i := 0; i < 10; i++ {
		urls[i], previewUrls[i] = uploadImage(t, "/app/test_files/sample.jpeg")
	}

	err := testImageClient.DeleteImagesFromS3(ctx, DeleteImagesFromS3Params{
		Urls: urls,
		PreviewUrls: previewUrls,
	})
	require.NoError(t, err)

	for _, url := range urls {
		imageBytes, err := testImageClient.GetImageFromS3(ctx, url)
		require.Error(t, err)
		require.Empty(t, imageBytes)

		previewImageBytes, err := testImageClient.GetPreviewFromS3(ctx, url)
		require.Error(t, err)
		require.Empty(t, previewImageBytes)
	}
}
