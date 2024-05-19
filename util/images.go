package util

import (
	"bytes"
	"io"
	"mime/multipart"
	// "net/textproto"
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

// package util

// import (
// 	"bytes"
// 	"io"
// 	"mime/multipart"
// 	"net/textproto"
// 	"os"
// 	// "path/filepath"
	
// 	"fmt"
// )

// func ImageAsFileHeader(path string) (*multipart.FileHeader, error) {
// 	// Open the file
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	// Create a buffer to store our multipart form
// 	body := &bytes.Buffer{}

// 	// Create a multipart writer
// 	writer := multipart.NewWriter(body)

// 	// Create a textproto.MIMEHeader
// 	h := make(textproto.MIMEHeader)
// 	h.Set("Content-Disposition", `form-data; name="image"; filename="`+ path + `"`)
// 	h.Set("Content-Type", "application/octet-stream")

// 	// Create the form file
// 	part, err := writer.CreatePart(h)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Write the file data to the part
// 	_, err = io.Copy(part, file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Close the multipart writer
// 	err = writer.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a *multipart.FileHeader for the form file
// 	fileHeader := multipart.FileHeader{
// 		Filename: path,
// 		Header:   h,
// 		Size:     int64(body.Len()),
// 	}
// 	fmt.Println("\n\n\nFile Name: ", fileHeader.Filename)

// 	return &fileHeader, nil
// }
