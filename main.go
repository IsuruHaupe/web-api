package main

import (
	"database/sql"
	"log"

	"github.com/IsuruHaupe/web-api/api"
	"github.com/IsuruHaupe/web-api/postgres"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/web-api?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	// Open the connection with the database.
	connection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database : ", err)
	}

	// Create a new PostgresConnection.
	postgresDatabase := postgres.NewPostgresConnection(connection)
	// Create the server.
	server := api.NewServer(postgresDatabase)

	// Start the server.
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Error when starting the server : ", err)
	}
}
