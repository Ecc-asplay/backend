package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Ecc-asplay/backend/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testQueries Querier
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error

	// Data
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot connect to config file: ", err)
	}

	// DB
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = NewStore(connPool)
	os.Exit(m.Run())
}
