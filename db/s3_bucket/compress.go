package db

import (
	"image/jpeg"
	"image/png"
	"io"
	compression "github.com/nurlantulemisov/imagecompression"
)

func CompressImage(image *multipart.File) (*multipart.File, error) {
	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf(err.Error())
	}

	compressing, err := compression.New(90)
	if err != nil {
		return nil, err
	}

	return compressing.Compress(img)
}
