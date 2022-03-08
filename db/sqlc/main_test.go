package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/web-api?sslmode=disable"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database : ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
