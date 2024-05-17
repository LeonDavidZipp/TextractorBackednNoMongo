package db

import (
	"context"
	"io"
	"os"
	"strings"
	"net/url"
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

func (c *Client) UploadImage(ctx context.Context, imageData []byte) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(uuid.New().String()),
		Body:   bytes.NewReader(imageData),
	}
	
	_, err := c.client.PutObject(
		ctx,
		input,
	)
	if err != nil {
		return "", err
	}
	
	location := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", *input.Bucket, *input.Key)

	return location, nil
}

func (c *Client) GetImage(ctx context.Context, link string) ([]byte, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	key := strings.TrimPrefix(parsed.Path, "/")

	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
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

	_, err := c.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Delete: &types.Delete{Objects: objectIds},
		},
	)
	return err
}
