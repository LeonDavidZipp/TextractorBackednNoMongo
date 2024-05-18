package util

import (
	"os"
	"mime/multipart"
	"net/textproto"
)


func ImageAsFileHeader(filePath string) (*multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileHeader := &multipart.FileHeader{
		Filename: filePath,
		Size: 0,
		Header: textproto.MIMEHeader{},
	}

	return fileHeader, nil
}
