package database

import db "github.com/IsuruHaupe/web-api/db/sqlc"

// We use embedding to ask for struct implementing this interface to implement every function in db.Querier.
// This is useful for switching between database implementation.
type Database interface {
	db.Querier
}
