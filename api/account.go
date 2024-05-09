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
	Email string ``
}

// Get Account
type getAccountParams struct {
	
}

// Delete Account


// MongoDB
// List Images