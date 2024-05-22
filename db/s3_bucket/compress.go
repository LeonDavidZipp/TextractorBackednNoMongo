package db

import (
	"image/jpeg"
	imglib "image"
	"mime/multipart"
	"io"
	"bytes"


	// compression "github.com/nurlantulemisov/imagecompression"
)

func CompressImage(image *multipart.File) (io.Reader, error) {
	// decode image
	img, _, err := imglib.Decode(*image)
	if err != nil {
		return nil, err
	}

	// compress image
	var buf bytes.Buffer
	var opts jpeg.Options
	opts.Quality = 75 // min 0 max 100
	err = jpeg.Encode(&buf, img, &opts)
	if err != nil {
		return nil, err
	}

	// turn image back to formfile
	reader := io.MultiReader(&buf)
	file := io.NopCloser(reader)

	// return
	return file, nil
}
