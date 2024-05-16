package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
)


type UploadImageParams struct {
	ID string `json:"id"`
	Image []byte `json:"image"`
}

// https://stackoverflow.com/questions/56744834/how-do-i-get-file-url-after-put-file-to-amazon-s3-in-go
func (u *S3Uploader) UploadImage(ctx context.Context, arg UploadImageParams) (string, error) {
	result, err := u.Uploader.Upload(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
			Key:    aws.String(arg.ID),
			Body:   bytes.NewReader(arg.Image),
		},
	)
	if err != nil {
		return "", err
	}
	return result.Location, err
}

func (u *S3Uploader) DeleteImage(ctx context.Context, location string) error {
	_, err := u.Uploader.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(location),
	})
	return err
}
