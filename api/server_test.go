package api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T) *Server {
	config := util.Config{
		FrontAddress:        []string{"*"},
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Error().Err(err).Msg("app.env が見つかりません")
		os.Exit(1)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Error().Err(err).Msg("データベースに接続できません")
		os.Exit(1)
	}
	t.Cleanup(func() {
		conn.Close()
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       1,
	})

	store := db.NewStore(conn)

	server, err := SetupRouter(config, store, rdb)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}
