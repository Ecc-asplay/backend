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
	// log 設定
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 環境変数の読み込み
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Error().Err(err).Msg("app.env が見つかりません")
		os.Exit(1)
	}

	// PostgreSQL 接続
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Error().Err(err).Msg("データベースに接続できません")
		os.Exit(1)
	}

	// マイグレーション実行
	initMigration(config.MigrationURL, config.DBSource)

	// DB 起動
	store := db.NewStore(conn)

	// Redis キャッシュ接続
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       1,
	})

	// サーバ設定
	server, err := api.SetupRouter(config, store, rdb)
	if err != nil {
		log.Error().Err(err).Msg("サーバを作成できません")
		os.Exit(1)
	}

	// サーバ起動
	log.Info().Msgf("Ginサーバに接続中: %s", config.HTTPServerAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("サーバを起動できません")
		os.Exit(1)
	}
}

func initMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("マイグレーションインスタンスの作成に失敗しました")
		os.Exit(1)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("マイグレーションの実行に失敗しました")
		os.Exit(1)
	}

	log.Info().Msg("データベースのマイグレーションが成功しました")
}
