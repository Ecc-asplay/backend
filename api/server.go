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

	// メモ：命名が気になるCreateuser→CreateUser
	// ユーザー関連のルート定義
	r.POST("/users", server.Createuser)
	r.DELETE("/users/:id", server.Deleteuser)
	r.GET("/users/:id", server.GetUserData)
	r.PUT("/users/:id/password", server.ResetPassword)
	r.PUT("/users/:id/disease-condition", server.UpdateDiseaseAndCondition)
	r.PUT("/users/:id/email", server.UpdateEmail)
	r.PUT("/users/:id/privacy", server.UpdateIsPrivacy)
	r.PUT("/users/:id/name", server.UpdateName)
	r.POST("/login", server.LoginUser)

	server.router = r

	return *server, nil

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
