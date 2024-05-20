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

	// Users (Postgres)
	router.POST("/users", server.createUser)

	router.GET("/users/:id", server.getUser)
	router.DELETE("/users/:id", server.deleteUser)
	
	// Images (Mongo)
	router.POST("/users/images", server.insertImage)
	router.GET("users/images", server.listImages)
	router.DELETE("/users/images", server.deleteImages)

	router.GET("/users/images/:id", server.getImage)
	router.PATCH("/users/images/:id", server.updateImage)

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
