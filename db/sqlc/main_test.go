package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/IsuruHaupe/web-api/config"
	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("Error when loading configuration in tests : ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database : ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
