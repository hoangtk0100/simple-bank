package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hoangtk0100/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
