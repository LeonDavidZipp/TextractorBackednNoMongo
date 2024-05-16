package db

import (
	"context"
)


// just first structure
type S3Client interface {
	UploadImage(ctx context.Context, arg UploadImageParams) (string, error)
	DeleteImage()
	DeleteImages()
}

func NewS3(cli S3Client) *Client {
	return &Client{
		client: cli,
	}
}

type Client struct {
	client S3Client
}
