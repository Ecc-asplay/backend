package main

import (
	"context"
	"os"

	"github.com/Ecc-asplay/backend/api"
	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config, err := util.LoadConfig("./")
	if err != nil {
		log.Info().Msg("app.env cannot find")
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Info().Msg("cannot connect to db")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       1,
	})

	store := db.NewStore(conn)
	server, err := api.SetupRouter(config, store, rdb)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
	}

	log.Info().Msgf("Connecting to Gin Server at %s", config.HTTPServerAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("cannot start server")
	}

}
