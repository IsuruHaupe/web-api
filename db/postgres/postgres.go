package postgres

import (
	"database/sql"

	"github.com/IsuruHaupe/web-api/db/database"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
)

// Postgres implementation, we use embedding to implement the database.Database interface for genericity purpose and make testing easier.
type PostgresDatabase struct {
	connection *sql.DB
	*db.Queries
}

// NewPostgresConnection creates a new PostgresDatabase object with a connection and a
// handler for calling database functions.
func NewPostgresConnection(connection *sql.DB) database.Database {
	return &PostgresDatabase{
		connection: connection,
		Queries:    db.New(connection),
	}
}
