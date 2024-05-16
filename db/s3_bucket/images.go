package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
)


type UploadImageParams struct {
	Image []byte
}

func (u *S3Uploader) UploadImage(ctx context.Context, arg UploadImageParams) error {
	_, err := u.Uploader.Upload(&s3.PutObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(arg.Key),
		Body:   bytes.NewReader(arg.Image),
	})
	return err
}
