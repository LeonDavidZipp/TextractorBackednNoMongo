package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
)

// Postgres
// Create Account
type createAccountParams struct {
	Owner string `json:"owner" binding:"reuired"`
	Email string `json:"email" binding:"required"`
}

// Get Account
type getAccountParams struct {
	ID int64 `json:"id" binding:"required"`
}

// Delete Account


// MongoDB
// List Images