package main

import (
	"context"
	"os"

	"github.com/go-redis/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Ecc-asplay/backend/api"
	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

func main() {
	// log 設置
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// env Data 取る
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Info().Msg("app.env cannot find")
	}
	// psql 接続
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Info().Msg("cannot connect to db")
	}

	// migration 実行
	initMigration(config.MigrationURL, config.DBSource)
	// DB 起動
	store := db.NewStore(conn)

	// redis Options settings
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       1,
	})

	// Server 設置
	server, err := api.SetupRouter(config, store, rdb)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
	}

	// server 起動
	log.Info().Msgf("Connecting to Gin Server at %s", config.HTTPServerAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("cannot start server")
	}
}

func initMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}
