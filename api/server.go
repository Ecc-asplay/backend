package api

import (
	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Querier
	router *gin.Engine
	config util.Config
}

func SetupRouter(config *util.Config, store db.Querier) (Server, error) {
	server := &Server{
		store:  store,
		config: *config,
	}
	r := gin.Default()

	r.GET("/", server.Createuser)

	server.router = r

	return *server, nil

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
