package db

import (

)


// just first structure
type S3Client interface {
	UploadImage(ctx context.Context, arg UploadImageParams) (string, error)
	DeleteImage()
	DeleteImages()
}
