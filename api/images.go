package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Upload Image
type uploadImageParams {

}

// Find Image
type findImageParams {
	ID int64 `json:"id" binding:"required"`
}

// Delete Image
type deleteImageParams {
	ID int64 `json:"id" binding:"required"`
}

// Update Image
type updateImageParams {
	ID          int64  `json:"id" binding:"required"`
	UpdatedText string `json:"updated_text" binding:"required"`
}

// List Images
type listImagesParams {
	AccountID int64 `json:"account_id" binding:"required"`
	Limit     int32 `json:"limit" binding:"required"`
	Offset    int32 `json:"offset" binding:"required"`
}

// Delete Images
type deleteImagesParams {
	ImageIDs []primitive.ObjectID `json:"image_ids" binding:"required"`
}