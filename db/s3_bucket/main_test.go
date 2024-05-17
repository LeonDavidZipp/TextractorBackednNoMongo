package db

import (
	"context"
	"testing"
	"os"
	"log"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
)

var testImageClient *Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Cannot load AWS config:", err)
	}

	s3Client := s3.NewFromConfig(config)
	testImageClient := NewS3(s3Client)

	os.Exit(m.Run())
}
