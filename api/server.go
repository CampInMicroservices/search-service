package api

import (
	"search-service/config"
	"search-service/db"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config config.Config
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config config.Config, store *db.Store) (*Server, error) {

	gin.SetMode(config.GinMode)
	router := gin.Default()

	server := &Server{
		config: config,
		store:  store,
	}

	// Setup routing for server
	v1 := router.Group("v1")
	{
		v1.GET("/listings/:id", server.GetListingByID)
		v1.GET("/listings", server.GetAllListings)
		v1.POST("/listings", server.CreateListing)
	}

	// Setup health check routes
	health := router.Group("health")
	{
		health.GET("/live", server.Live)
		health.GET("/ready", server.Ready)
	}

	// TODO: Setup metrics routes

	server.router = router
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
