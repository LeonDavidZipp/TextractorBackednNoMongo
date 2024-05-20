package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
)

// Postgres
// Create User
type createUserRequest struct {
	Owner      string         `json:"owner" binding:"required"`
	Email      string         `json:"email" binding:"required"`
	GoogleID   *string        `json:"google_id"`
	FacebookID *string        `json:"facebook_id"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var googleID, facebookID sql.NullString

	if req.GoogleID != nil {
		googleID = sql.NullString{
			String: *req.GoogleID,
			Valid:  true,
		}
	}
	if req.FacebookID != nil {
		facebookID = sql.NullString{
			String: *req.FacebookID,
			Valid:  true,
		}
	}

	arg := db.CreateUserParams{
		Owner:      req.Owner,
		Email:      req.Email,
		GoogleID:   googleID,
		FacebookID: facebookID,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Get User
type getUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Delete User
type deleteUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteUser(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
