package db

import (
	"os"
	"io"
	"context"
	"testing"
	"github.com/stretchr/testify/require"
)


func imageToBytes(imagePath string) ([]byte, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func uploadImage(t *testing.T, imagePath string) string {
	ctx := context.Background()

	imageBytes, err := imageToBytes(imagePath)
	require.NoError(t, err)

	result, err := testImageClient.UploadAndExtractImage(ctx, imageBytes)
	link := result.Link
	text := result.Text

	require.NoError(t, err)
	require.NotEmpty(t, link)
	require.NotEmpty(t, text)

	return link
}

func TestUploadImage(t *testing.T) {
	uploadImage(t, "../../test_files/sample.jpeg")
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()

	link := uploadImage(t, "../../test_files/sample.jpeg")

	imageBytes, err := testImageClient.GetImage(ctx, link)
	require.NoError(t, err)
	require.NotEmpty(t, imageBytes)
}

func TestDeleteImages(t *testing.T) {
	ctx := context.Background()

	links := make([]string, 10)
	for i := 0; i < 10; i++ {
		links[i] = uploadImage(t, "../../test_files/sample.jpeg")
	}

	err := testImageClient.DeleteImagesFromS3(ctx, links)
	require.NoError(t, err)

	for _, link := range links {
		imageBytes, err := testImageClient.GetImage(ctx, link)
		require.Error(t, err)
		require.Empty(t, imageBytes)
	}
}
