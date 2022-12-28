package main

import (
	"log"
	"search-service/api"
	"search-service/config"
	"search-service/db"
)

func main() {

	// Load configuration settings
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Connect to the database
	store, err := db.Connect(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Run DB migration
	db.RunDBMigration(config.MigrationURL, config.DBSource)

	// Create a server and setup routes
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create a server: ", err)
	}

	// Start a server
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("Failed to start a server: ", err)
	}
}
