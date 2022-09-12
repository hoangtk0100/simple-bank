package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hoangtk0100/simple-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
