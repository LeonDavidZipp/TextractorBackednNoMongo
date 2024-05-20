package api

import (
	"database/sql"
	"net/http"
	"errors"
	"mime/multipart"
	"github.com/gin-gonic/gin"
	st "github.com/LeonDavidZipp/Textractor/db/store"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
)


type insertImageRequest struct {
	UserID int64                 `form:"user_id" binding:"required"`
	Image  *multipart.FileHeader `form:"image" binding:"required"`
}

func (s *Server) insertImage(ctx *gin.Context) {
	var req insertImageRequest
	
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := st.UploadImageTransactionParams{
		UserID: req.UserID,
		Image: req.Image,
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
	ID int64 `json:"id" binding:"required"`
}

func (s *Server) getImage(ctx *gin.Context) {
	var req findImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	image, err := s.store.GetImageFromSQL(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, image)

}

// Update Image
type updateImageRequest struct {
	ID          int64  `json:"id" binding:"required"`
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
	UserID int64 `json:"user_id" binding:"required"`
	Limit  int64 `json:"limit" binding:"required"`
	Offset int64 `json:"offset" binding:"required"`
}

func (s *Server) listImages(ctx *gin.Context) {
	var req listImagesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	arg := db.ListImagesParams{
		UserID: req.UserID,
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
	ImageIDs []int64  `json:"image_ids" binding:"required"`
	Urls	 []string `json:"urls" binding:"required"`
}

func (s *Server) deleteImages(ctx *gin.Context) {
	var req deleteImagesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteImageTransaction(ctx, req.ImageIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
