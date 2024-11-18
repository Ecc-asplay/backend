package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Ecc-asplay/backend/util"
)

var (
	testQueries Store
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       1,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	os.Exit(m.Run())
}
