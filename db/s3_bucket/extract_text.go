package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

func ExtractText(ctx context.Context, link string) (string, error) {
	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	client := textract.NewFromConfig(config)
	_, err := client.DetectDocumentText(ctx, &textract.DetectDocumentTextInput{
		Document: &textract.Document{
			S3Object: &textract.S3Object{
				Bucket: aws.String("AWS_BUCKET_NAME"),
				Name:   aws.String(KeyFromLink(link),
				),
			},
		},
	})
	if err != nil {
		return "", err
	}
}
