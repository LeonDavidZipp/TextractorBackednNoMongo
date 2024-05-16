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
)


type UploadImageParams struct {
	// ID string `json:"id"`
	Image []byte `json:"image"`
}

// https://stackoverflow.com/questions/56744834/how-do-i-get-file-url-after-put-file-to-amazon-s3-in-go
func (c *Client) UploadImage(ctx context.Context, arg UploadImageParams) (string, error) {
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

func (c *Client) DeleteImages(ctx context.Context, links []string) error {
	for _, link := range links {
		err := c.DeleteImage(ctx, link)
		if err != nil {
			return err
		}
	}
	return nil
}
