package postgres

import (
	"database/sql"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
)

type Database interface {
	db.Querier
}

type PostgresDatabase struct {
	connection *sql.DB
	*db.Queries
}

func NewPostgresConnection(connection *sql.DB) Database {
	return &PostgresDatabase{
		connection: connection,
		Queries:    db.New(connection),
	}
}
