package api

import (
	conf "search-service/config"
	"search-service/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config conf.Config
	store  *db.Store
	router *gin.Engine
}

func NewServer(config conf.Config, store *db.Store) (*Server, error) {

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

	go conf.WatchConsulConfig("DB_SOURCE", server.config.ConsulAddress, func(source string) {
		store, err := db.Connect(server.config.DBDriver, source)

		if err == nil {
			// Run DB migration
			db.RunDBMigration(server.config.MigrationURL, server.config.DBSource)

			// Rewire the connection
			server.store.Close()
			server.store = store
		}
	})

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
