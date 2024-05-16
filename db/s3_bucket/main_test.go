package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
)

var testImageClient *Client

func TestMain(m *testing.M) {
	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Cannot load AWS config:", err)
	}

	s3Client := s3.NewFromConfig(config)
	testImageClient := NewS3(s3Client)

	os.exit(m.Run())
}
