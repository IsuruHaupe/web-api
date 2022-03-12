package postgres

import (
	"database/sql"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
)

// We use embedding to ask for struct implementing this interface to implement every function in db.Querier.
// This is useful for switching between database implementation.
type Database interface {
	db.Querier
}

// Postgres implementation.
type PostgresDatabase struct {
	connection *sql.DB
	*db.Queries
}

// NewPostgresConnection creates a new PostgresDatabase object with a connection and a
// handler for calling database functions.
func NewPostgresConnection(connection *sql.DB) Database {
	return &PostgresDatabase{
		connection: connection,
		Queries:    db.New(connection),
	}
}
