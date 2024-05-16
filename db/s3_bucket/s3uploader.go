package db

import (

)


// just first structure
type S3Uploader interface {
	UploadImage()
	DeleteImage()
	DeleteImages()
}
