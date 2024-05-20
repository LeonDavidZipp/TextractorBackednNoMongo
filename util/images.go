package util

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func ImageAsFileHeader(path string) (*multipart.FileHeader, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffer to store our multipart form
	body := &bytes.Buffer{}

	// Create a multipart writer
	writer := multipart.NewWriter(body)

	// Create a form file
	part, err := writer.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		return nil, err
	}

	// Write the file data to the part
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a *multipart.FileHeader for the form file
	fileHeader := multipart.FileHeader{
		Filename: filepath.Base(path),
		Header:   part.Header,
		Size:     int64(body.Len()),
	}

	return &fileHeader, nil
}
