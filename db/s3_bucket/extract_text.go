package db

import (
	"context"
	"strings"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

func ExtractText(ctx context.Context, link string) (string, error) {
	config, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return "", err
	}

	key, err := KeyFromLink(ctx, link)
	if err != nil {
		return "", err
	}

	client := textract.NewFromConfig(config)
	result, err := client.DetectDocumentText(ctx, &textract.DetectDocumentTextInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String("AWS_BUCKET_NAME"),
				Name:   aws.String(key),
			},
		},
	})
	if err != nil {
		return "", err
	}

	var documentText strings.Builder
	for _, block := range result.Blocks {
		if block.Text != nil {
			documentText.WriteString(*block.Text + " ")
		}
	}

	return documentText.String(), nil
}
