package db

import (
	"context"
	"net/url"
	"strings"
	"fmt"
)

func KeyFromLink(link string) (string, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	return strings.TrimPrefix(parsed.Path, "/"), nil
}

func LinkFromKey(key string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", os.Getenv("AWS_BUCKET_NAME"), key)
}
