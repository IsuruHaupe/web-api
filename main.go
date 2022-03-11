package main

import (
	"database/sql"
	"log"

	"github.com/IsuruHaupe/web-api/api"
	"github.com/IsuruHaupe/web-api/config"
	"github.com/IsuruHaupe/web-api/db/postgres"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration using viper
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Error when loading configuration : ", err)
	}

	// Open the connection with the database.
	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database : ", err)
	}

	// Create a new PostgresConnection.
	postgresDatabase := postgres.NewPostgresConnection(connection)
	// Create the server.
	server, err := api.NewServer(config, postgresDatabase)
	if err != nil {
		log.Fatal("cannot create server : ", err)
	}

	// Start the server.
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Error when starting the server : ", err)
	}
}
