package api

import (
	"github.com/gin-gonic/gin"
	st "github.com/LeonDavidZipp/Textractor/db/store"
)


type Server struct {
	store st.Store
	router *gin.Engine
}

func NewServer(store st.Store) *Server {
	server := &Server{
		store : store,
	}

	router := gin.Default()

	// Accounts (Postgres)
	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	
	// Images (Mongo)
	router.POST("/accounts/images", server.insertImage)
	router.GET("accounts/images", server.listImages)
	router.DELETE("/accounts/images", server.deleteImages)

	router.GET("/accounts/images/:id", server.findImage)
	router.DELETE("/accounts/images/:id", server.deleteImage)
	router.PATCH("/accounts/images/:id", server.updateImage)

	server.router = router

	return server
}

// Starts http server on specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// formats error response into json body
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
