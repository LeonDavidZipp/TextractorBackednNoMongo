package util

import (
	"os"
	"io"
	"bytes"
	"errors"
	"path/filepath"
	"mime"
	"mime/multipart"
)


func ImageAsFileHeader(filePath string) (*multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	writer := multipart.NewWriter(body)
	defer writer.Close()

	formWriter, err := writer.CreateFormFile("image", "image.jpg")

	_, err = io.Copy(formWriter, file)
	if err != nil {
		return {}, err
	}
}


