package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"context"
	"bytes"
	"os"
	"github.com/google/uuid"
	"io"
)


// https://stackoverflow.com/questions/56744834/how-do-i-get-file-url-after-put-file-to-amazon-s3-in-go
func (c *Client) UploadImage(ctx context.Context, imageData []bytes) (string, error) {
	uploader := manager.NewUploader(c.client)

	result, err := uploader.Upload(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
			Key:    aws.String(uuid.New().String()),
			Body:   bytes.NewReader(arg.Image),
		},
	)
	if err != nil {
		return "", err
	}
	return result.Location, err
}

func (c *Client) GetImage(ctx context.Context, link string) ([]byte, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	key := strings.TrimPrefix(parsed.Path, "/")

	result, err := basics.S3Client.GetObject(ctx, &s3.GetObjectInput{
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

func (c *Client) DeleteImages(ctx context.Context, links []string) error {
	var objectIds []types.ObjectIdentifier
	for _, link := range links {
		parsed, err := url.Parse(link)
		if err != nil {
			return err
		}
		key := strings.TrimPrefix(parsed.Path, "/")
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err = c.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Delete: &types.Delete{Objects: objectIds},
		},
	)
	return err
}
