package db

import (
	"context"
	"io"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"mime/multipart"

	"fmt"
)

type UploadImageResult struct {
	URL string `json:"url"`
	Text string `json:"text"`
}

func (c *Client) UploadAndExtractImage(ctx context.Context, image *multipart.FileHeader) (UploadImageResult, error) {
	img, err := image.Open()
	if err != nil {
		return UploadImageResult{}, err
	}
	defer img.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(fmt.Sprintf("%s_orig_name_after_%s", uuid.New().String(), image.Filename)),
		Body:   img,
	}

	uploadResult, err := c.Uploader.Upload(
		ctx,
		input,
	)
	if err != nil {
		return UploadImageResult{}, err
	}

	url := uploadResult.Location

	text, err := ExtractText(ctx, *input.Key)
	if err != nil {
		return UploadImageResult{}, err
	}
	
	result := UploadImageResult{
		URL: url,
		Text: text,
	}

	return result, nil
}

func (c *Client) GetImage(ctx context.Context, url string) ([]byte, error) {
	key, err := KeyFromURL(ctx, url)
	if err != nil {
		return nil, err
	}

	result, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	imageData, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func (c *Client) DeleteImagesFromS3(ctx context.Context, urls []string) error {
	var objectIds []types.ObjectIdentifier
	for _, url := range urls {
		key, err := KeyFromURL(ctx, url)
		if err != nil {
			return err
		}

		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err := c.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Delete: &types.Delete{Objects: objectIds},
		},
	)
	return err
}
