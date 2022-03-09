package postgres

import (
	"database/sql"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
)

type PostgresDatabase struct {
	connection *sql.DB
	*db.Queries
}

func NewPostgresConnection(connection *sql.DB) *PostgresDatabase {
	return &PostgresDatabase{
		connection: connection,
		Queries:    db.New(connection),
	}
}
