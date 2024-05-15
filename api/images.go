package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert Image
type insertImageRequest {
	AccountID int64  `json:"account_id" binding:"required"`
	Text      string `bson:"text" json:"text"`
	// link to the image in s3 storage
	Link      string `bson:"link" json:"link"`
	Image64   string `bson:"image_64" json:"image_64"`
}

func (s *Server) insertImage(ctx *gin.Context) {
	var req insertImageRequest

	if err := ctx.ShouldBindJson(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg :=mongodb.InsertImageParams{
		AccountID: req.AccountID,
		Text: req.Text,
		Link: req.Link,
		Image64: req.Image64,
	}

	image, err := s.store.InsertImage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, image)
}

// Find Image
type findImageRequest {
	ID primitive.ObjectID `json:"id" binding:"required"`
}

func (s *Server) findImage(ctx *gin.Context) {
	var req findImageRequest

	if err := ctx.ShouldBindJson(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	image, err := s.store.FindImage(ctx, req.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, image)

}

// Delete Image
type deleteImageRequest {
	ID primitive.ObjectID `json:"id" binding:"required"`
}

func (s *Server) deleteImage(ctx *gin.Context) {
	var req deleteImageRequest

	if err := ctx.ShouldBindJson(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

}

// Update Image
type updateImageRequest {
	ID          primitive.ObjectID  `json:"id" binding:"required"`
	UpdatedText string `json:"updated_text" binding:"required"`
}

func (s *Server) updateImage(ctx *gin.Context) {
	var req updateImageRequest

	if err := ctx.ShouldBindJson(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

}

// List Images
type listImagesRequest {
	AccountID int64 `json:"account_id" binding:"required"`
	Limit     int32 `json:"limit" binding:"required"`
	Offset    int32 `json:"offset" binding:"required"`
}

func (s *Server) listImage(ctx *gin.Context) {
	var req listImageRequest

	if err := ctx.ShouldBindJson(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

}

// Delete Images
type deleteImagesRequest {
	ImageIDs []primitive.ObjectID `json:"image_ids" binding:"required"`
}

func (s *Server) deleteImages(ctx *gin.Context) {
	var req deleteImagesRequest

	if err := ctx.ShouldBindJson(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

}
