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
	Link string `json:"link"`
	Text string `json:"text"`
}

func (c *Client) UploadAndExtractImage(ctx context.Context, image *multipart.FileHeader) (UploadImageResult, error) {
	img, err := image.Open()
	if err != nil {
		fmt.Println("\n\n\nError opening image: ", err)
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

	link := uploadResult.Location
	fmt.Println("\n\n\nLink: ", link)

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
