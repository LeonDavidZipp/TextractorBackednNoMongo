package db

import (
	"context"
	"strings"
	"fmt"
	"os"
	"net/url"
)

func KeyFromUrl(ctx context.Context, u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(parsed.Path, "/"), nil
}

func UrlFromKey(ctx context.Context, key string) string {
	return fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET_NAME"), key)
}
