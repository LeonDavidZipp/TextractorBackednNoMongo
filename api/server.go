package api

// import (
// 	"database/sql"
// 	"net/http"
// 	"github.com/gin-gonic/gin"
// 	db "github.com/LeonDavidZipp/Textractor/db"
// 	sqldb "github.com/LeonDavidZipp/Textractor/db/sqlc"
// 	mongodb "github.com/LeonDavidZipp/Textractor/db/mongo_db"
// )


// type Server struct {
// 	store db.Store
// 	router *gin.Engine
// }

// func NewServer(store db.Store) *Server {
// 	server := &Server{
// 		store : store,
// 	}

// 	router := gin.Default()

// 	// Accounts (Postgres)
// 	router.POST("/accounts", server.createAccount)
// 	router.GET("/accounts/:id", server.getAccount)
// 	router.GET("accounts", server.listAccounts)
// 	router.DELETE("/accounts/:id", server.deleteAccount)
	
// 	// Images (Mongo)
// 	router.POST("/accounts", server.uploadAccount)
// 	router.GET("/accounts/images/:id", server.getAccount)
// 	router.GET("accounts/images", server.listAccounts)
// 	router.DELETE("/accounts/images/:id", server.deleteImage)

// 	server.router = router

// 	return server
// }

// // Starts http server on specified address
// func (server *Server) Start(address string) error {
// 	return server.router.Run(address)
// }

// // formats error response into json body
// func errorResponse(err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }
