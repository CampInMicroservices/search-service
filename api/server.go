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

var requestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "rso_http_request_counter",
		Help: "Number of HTTP requests",
	})

var notFoundCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "rso_not_found_counter",
		Help: "Number of HTTP not found responses",
	})

var lastRequestTime = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "rso_last_request_received_time",
		Help: "Time when the last request was processed",
	})

// Server definition
func NewServer(config conf.Config, store *db.Store) (*Server, error) {

	gin.SetMode(config.GinMode)
	router := gin.Default()
	router.Use(middleware.Logger(config.LogitAddress))

	// Prometheus gauge registration
	registerCounter(requestCounter)
	registerCounter(notFoundCounter)
	registerGauge(lastRequestTime)

	// Metrics middleware for all requests
	router.Use(func(context *gin.Context) {
		if context.Request.URL.Path == "/metrics" {
			return
		}

		requestCounter.Inc()
		lastRequestTime.SetToCurrentTime()

		// Process request
		context.Next()

		if context.Writer.Status() == 404 {
			notFoundCounter.Inc()
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

func registerCounter(c prometheus.Counter) {
	err := prometheus.Register(c)
	handleErr(err)
}

func registerGauge(g prometheus.Gauge) {
	err := prometheus.Register(g)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
