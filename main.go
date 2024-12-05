package main

import (
	"context"
	"os"

	"github.com/go-redis/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Ecc-asplay/backend/api"
	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
	"github.com/Ecc-asplay/backend/worker"
)

func main() {
	// log 設置
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// env Data 取る
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Error().Err(err).Msg("app.env cannot find")
		os.Exit(1)
	}

	// psql 接続
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Error().Err(err).Msg("cannot connect to db")
		os.Exit(1)
	}

	// migration 実行
	initMigration(config.MigrationURL, config.DBSource)

	// DB 起動
	store := db.NewStore(conn)

	// redis Options settings
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
		DB:   0,
	}
	// redis キャッシュ 接続
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       1,
	})

	// Processer
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	runTaskProcessor(redisOpt, store)

	// Server 設置
	server, err := api.SetupRouter(config, store, rdb, taskDistributor)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
		os.Exit(1)
	}

	// server 起動
	log.Info().Msgf("Connecting to Gin Server at %s", config.HTTPServerAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("cannot start server")
		os.Exit(1)
	}
}

func initMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
		os.Exit(1)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
		os.Exit(1)
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
