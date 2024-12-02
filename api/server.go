package api

import (
	"fmt"
	"io"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
	"github.com/Ecc-asplay/backend/worker"
)

type Server struct {
	store           db.Store
	router          *gin.Engine
	redis           *redis.Client
	config          util.Config
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func SetupRouter(config util.Config, store db.Store, redis *redis.Client, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:           store,
		config:          config,
		redis:           redis,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	server.GinRequest(config)

	return server, nil
}

func (server *Server) GinRequest(config util.Config) {
	// Log file cancel
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	// Gin Start
	r := gin.Default()
	r.Use(GinLogger())
	corsConfig := cors.Config{
		AllowOrigins:     config.FrontAddress,
		AllowMethods:     config.AllowHeaders,
		AllowHeaders:     []string{"Authorization", "TokenID", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           config.AccessTokenDuration,
	}

	r.Use(cors.New(corsConfig))

	r.POST("/users", server.CreateUser)
	r.POST("/login", server.LoginUser)
	r.GET("/getposts", server.GetAllPost)

	authRoutes := r.Group("/").Use(authMiddleware(server.tokenMaker))

	// User
	r.DELETE("/users/:id", server.DeleteUser)
	r.GET("/users/:id", server.GetUserData)
	r.PUT("/users/:id/password", server.ResetPassword)
	r.PUT("/users/:id/disease-condition", server.UpdateDiseaseAndCondition)
	r.PUT("/users/:id/email", server.UpdateEmail)
	r.PUT("/users/:id/privacy", server.UpdateIsPrivacy)
	r.PUT("/users/:id/name", server.UpdateName)

	// Posts
	r.POST("/createpost", server.CreatePost)
	r.DELETE("/delpost", server.DeletePost)

	// Tag
	r.POST("/tag/add", server.CreateTag)
	r.POST("/tag/get", server.GetTag)

	// bookmark
	authRoutes.POST("createbookmark", server.CreateBookmark)
	authRoutes.DELETE("deletebookmark", server.DeleteBookmark)
	r.GET("getbookmark", server.GetBookmark)

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
