package api

import (
	"net/http"
	"search-service/config"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config config.Config
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config config.Config) (*Server, error) {
	server := &Server{
		config: config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/test", server.getTest)

	server.router = router
}

func (server *Server) getTest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
