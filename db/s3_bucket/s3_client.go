package db

import (
	"context"
)


// just first structure
type S3Client interface {
	UploadImage(ctx context.Context, arg UploadImageParams) (string, error)
	GetImage(ctx context.Context, link string) ([]byte, error)
	DeleteImages(ctx context.Context, links []string) error
}

func NewS3(cli S3Client) *Client {
	return &Client{
		client: cli,
	}
}

type Client struct {
	client S3Client
}
