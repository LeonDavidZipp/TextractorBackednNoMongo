package api

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	db "github.com/LeonDavidZipp/Textractor/db/sqlc"
)


type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server { // testing: remove later
	server := &Server{
		store : store,
	}

	router := gin.Default()
	// Accounts (Postgres)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("accounts", server.listAccounts)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	
	// Images (Mongo)
	// router.POST("/accounts", server.uploadAccount)
	// router.GET("/accounts/images/:id", server.getAccount)
	// router.GET("accounts/images", server.listAccounts)
	// router.DELETE("/accounts/images/:id", server.deleteImage)

	server.router = router

	return server
}
