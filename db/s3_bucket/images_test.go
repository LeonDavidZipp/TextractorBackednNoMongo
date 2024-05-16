package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"bytes"
	"os"
	"io"
	"context"
	"time"
	"testing"
	"github.com/stretchr/testify/require"
)


func TestUploadImage(t *testing.T) {
	ctx := context.Background()
}

func TestGetImage(t *testing.T) {
	ctx := context.Background()
}

func TestDeleteImages(t *testing.T) {
	ctx := context.Background()
}
