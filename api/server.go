package api

import (
	"fmt"
	"io"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	store      db.Querier
	redis      *redis.Client
	router     *gin.Engine
	config     util.Config
	tokenMaker token.Maker
}

func SetupRouter(config *util.Config, store db.Querier, rdb *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		redis:      rdb,
		config:     *config,
		tokenMaker: tokenMaker,
	}

	// Set Log
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.Use(GinLogger())

	r.GET("/", server.Createuser)

	server.router = r

	return server, nil

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
