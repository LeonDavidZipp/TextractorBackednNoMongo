package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)


// just first structure
// type S3Client interface {
// 	UploadImage(ctx context.Context, imageData []byte) (string, error)
// 	GetImage(ctx context.Context, link string) ([]byte, error)
// 	DeleteImages(ctx context.Context, links []string) error
// }

type S3Client interface {
	manager.UploadAPIClient
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func NewS3(cli S3Client) *Client {
	return &Client{S3Client: cli}
}

type Client struct {
	S3Client
}
