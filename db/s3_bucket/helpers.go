package db

import (
	"context"
	"net/url"
	"strings"
	"fmt"
	"os"
)

func KeyFromLink(ctx context.Context, link string) (string, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(parsed.Path, "/"), nil
}

func LinkFromKey(ctx context.Context, key string) string {
	return fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET_NAME"), key)
}
