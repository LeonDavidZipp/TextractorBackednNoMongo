package api

import (
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	st "github.com/LeonDavidZipp/Textractor/db/store"
)

type insertImageRequest struct {
	AccountID int64  `json:"account_id" binding:"required"`
	// Filepath  string `json:"filepath" binding:"required"`
	ImageData []byte `json:"image_data" binding:"required"`
}

// func (s *Server) insertImage(ctx *gin.Context) {
// 	var req insertImageRequest

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := st.UploadImageTransactionParams{
// 		AccountID: req.AccountID,
// 		ImageData: req.ImageData,
// 	}

// 	result, err := s.store.UploadImageTransaction(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, result)
// }

func (s *Server) insertImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accountID := ctx.PostForm("account_id")
	if accountID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	arg := st.UploadImageTransactionParams{
		AccountID: accountID,
		File: file,
	}

	result, err := s.store.UploadImageTransaction(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Find Image
type findImageRequest struct {
	ID primitive.ObjectID `json:"id" binding:"required"`
}

func (s *Server) findImage(ctx *gin.Context) {
	var req findImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
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
type deleteImageRequest struct {
	ID primitive.ObjectID `json:"id" binding:"required"`
}

func (s *Server) deleteImage(ctx *gin.Context) {
	var req deleteImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteImage(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)

}

// Update Image
type updateImageRequest struct {
	ID          primitive.ObjectID  `json:"id" binding:"required"`
	UpdatedText string `json:"updated_text" binding:"required"`
}

func (s *Server) updateImage(ctx *gin.Context) {
	var req updateImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

}

// List Images
type listImagesRequest struct {
	AccountID int64 `json:"account_id" binding:"required"`
	Limit     int64 `json:"limit" binding:"required"`
	Offset    int64 `json:"offset" binding:"required"`
}

func (s *Server) listImages(ctx *gin.Context) {
	var req listImagesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	arg := mongodb.ListImagesParams{
		AccountID: req.AccountID,
		Limit: req.Limit,
		Offset: req.Offset,
	}

	images, err := s.store.ListImages(ctx, arg)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, images)
}

// Delete Images
type deleteImagesRequest struct {
	ImageIDs []primitive.ObjectID `json:"image_ids" binding:"required"`
}

func (s *Server) deleteImages(ctx *gin.Context) {
	var req deleteImagesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteImagesFromMongo(ctx, req.ImageIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
