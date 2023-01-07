package api

import (
	"log"
	conf "search-service/config"
	"search-service/db"
	"search-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type Server struct {
	config conf.Config
	store  *db.Store
	router *gin.Engine
}

func NewServer(config conf.Config, store *db.Store) (*Server, error) {

	gin.SetMode(config.GinMode)
	router := gin.Default()
	router.Use(middleware.Logger(config.LogitAddress))

	// Prometheus gauge registration
	lastRequestReceivedTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "last_request_received_time",
		Help: "Time when the last request was processed",
	})
	err := prometheus.Register(lastRequestReceivedTime)
	if err != nil {
		log.Fatalln(err)
	}

	// Middleware to set lastRequestReceivedTime for all requests
	router.Use(func(context *gin.Context) {
		log.Println(context.Request.URL.Path)
		if context.Request.URL.Path != "/metrics" {
			lastRequestReceivedTime.SetToCurrentTime()
		}
	})

	server := &Server{
		config: config,
		store:  store,
	}

	// Server
	v1 := router.Group("v1")
	{
		v1.GET("/listings/:id", server.GetListingByID)
		v1.GET("/listings", server.GetAllListings)
		v1.POST("/listings", server.CreateListing)

		v1.GET("/recommendations", server.GetRecommendedLocations)
	}

	// Health
	health := router.Group("health")
	{
		health.GET("/live", server.Live)
		health.GET("/ready", server.Ready)
	}

	// Metrics
	metrics := router.Group("metrics")
	{
		metrics.GET("", server.Metrics)
	}

	server.router = router
	return server, nil
}

// Start HTTP server and initialize consul dynamic config
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
