package main

import (
	"log"
	"search-service/api"
	"search-service/config"
	"search-service/db"
)

//	@title			CampIn Search Service API
//	@version		1.0
//	@description	This is a search service server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Šimen Ravnik
//	@contact.email	sr8905@student.uni-lj.si

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		20.13.80.52
// @BasePath	search-service/v1
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
