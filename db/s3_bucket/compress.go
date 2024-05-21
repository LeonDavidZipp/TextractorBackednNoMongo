package db

import (
	"image/jpeg"
	imglib "image"
	"mime/multipart"

	compression "github.com/nurlantulemisov/imagecompression"
)

func CompressImage(image *multipart.File) (*multipart.File, error) {
	img, _, err := imglib.Decode(image)
	if err != nil {
		return nil, err
	}

	// Compress the image
	var buf bytes.Buffer
	var opts jpeg.Options
	opts.Quality = 75
	err = jpeg.Encode(&buf, img, &opts)
	if err != nil {
		return nil, err
	}

	
}
