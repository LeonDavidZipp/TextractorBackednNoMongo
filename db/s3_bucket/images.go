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
	Url        string `json:"url"`
	PreviewUrl string `json:"preview_url"`
	Text       string `json:"text"`
}

func (c *Client) UploadAndExtractImage(ctx context.Context, image *multipart.FileHeader) (UploadImageResult, error) {
	img, err := image.Open()
	if err != nil {
		return UploadImageResult{}, err
	}
	defer img.Close()

	// upload to image storage
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

	url := uploadResult.Location

	// extract text from image
	text, err := ExtractText(ctx, *input.Key)
	if err != nil {
		return UploadImageResult{}, err
	}

	// upload preview storage
	compressedImg, err := CompressImage(image)
	if err != nil {
		return UploadImageResult{}, err
	}

	previewInput := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_PREVIEW_BUCKET_NAME")),
		Key:    aws.String(fmt.Sprintf("%s_preview_%s", uuid.New().String(), image.Filename)),
		Body:   compressedImg,
	}

	previewUploadResult, err = c.Uploader.Upload(
		ctx,
		previewInput,
	)
	if err != nil {
		return UploadImageResult{}, err
	}

	previewUrl := previewUploadResult.Location
	
	// return result
	result := UploadImageResult{
		Url: url,
		PreviewUrl: previewUrl,
		Text: text,
	}

	return result, nil
}

func (c *Client) GetImageFromS3(ctx context.Context, url string) ([]byte, error) {
	key, err := KeyFromUrl(ctx, url)
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

type DeleteImagesFromS3Params struct {
	Urls        []string `json:"urls"`
	PreviewUrls []string `json:"preview_urls"`
}

func (c *Client) DeleteImagesFromS3(ctx context.Context, arg DeleteImagesFromS3Params) error {
	var objectIds []types.ObjectIdentifier
	var previewObjectIds []types.ObjectIdentifier

	for i, url := range arg.Urls {
		key, err := KeyFromUrl(ctx, url)
		if err != nil {
			return err
		}

		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})

		previewKey, err := KeyFromUrl(ctx, arg.PreviewUrls[i])
		if err != nil {
			return err
		}

		previewObjectIds = append(previewObjectIds, types.ObjectIdentifier{Key: aws.String(previewKey)})
	}

	_, err := c.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		return err
	}

	_, err = c.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_PREVIEW_BUCKET_NAME")),
		Delete: &types.Delete{Objects: previewObjectIds},
	})

	return err
}
