package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
)

// Postgres
// Create Account
type createAccountRequest struct {
	Owner      string         `json:"owner" binding:"reuired"`
	Email      string         `json:"email" binding:"required"`
	GoogleID   sql.NullString `json:"google_id"`
	FacebookID sql.NullString `json:"facebook_id"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Email: req.Email,
		GoogleID: req.GoogleID,
		FacebookID: req.FacebookID,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Get Account
type getAccountRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Delete Account
type deleteAccountRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (s *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
