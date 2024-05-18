package db

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/google/uuid"
	"mime/multipart"
)

type UploadImageResult struct {
	Link string `json:"link"`
	Text string `json:"text"`
}

// func (c *Client) UploadAndExtractImage(ctx context.Context, file *multipart.FileHeader) (UploadImageResult, error) {
// 	f, err := file.Open()
// 	if err != nil {
// 		return UploadImageResult{}, err
// 	}

// 	input := &s3.PutObjectInput{
// 		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
// 		Key:    aws.String(uuid.New().String()),
// 		Body:   f,
// 	}

// 	_, err := c.PutObject(
// 		ctx,
// 		input,
// 	)
// 	if err != nil {
// 		return UploadImageResult{}, err
// 	}
	
// 	link := LinkFromKey(ctx, *input.Key)

// 	text, err := ExtractText(ctx, link)
// 	if err != nil {
// 		return UploadImageResult{}, err
// 	}
	
// 	result := UploadImageResult{
// 		Link: link,
// 		Text: text,
// 	}

// 	return result, nil
// }

func (c *Client) UploadAndExtractImage(ctx context.Context, file *multipart.FileHeader) (UploadImageResult, error) {
	f, err := file.Open()
	if err != nil {
		return UploadImageResult{}, err
	}
	defer f.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(uuid.New().String()),
		Body:   f,
	}

	uploader := manager.NewUploader(c)

	uploadResult, err := uploader.Upload(
		ctx,
		input,
	)
	if err != nil {
		return UploadImageResult{}, err
	}

	link := uploadResult.Location

	text, err := ExtractText(ctx, link)
	if err != nil {
		return UploadImageResult{}, err
	}
	
	result := UploadImageResult{
		Link: link,
		Text: text,
	}

	return result, nil
}

func (c *Client) GetImage(ctx context.Context, link string) ([]byte, error) {
	key, err := KeyFromLink(ctx, link)
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

func (c *Client) DeleteImagesFromS3(ctx context.Context, links []string) error {
	var objectIds []types.ObjectIdentifier
	for _, link := range links {
		key, err := KeyFromLink(ctx, link)
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
