package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

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
		return nil, fmt.Errorf("トークンメーカーの作成に失敗しました: %w", err)
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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           config.AccessTokenDuration,
	}

	r.Use(cors.New(corsConfig))

	// ログイン前
	r.POST("/users", server.CreateUser)
	r.POST("/login", server.LoginUser)
	r.GET("/post/getall", server.GetAllPost)
	r.POST("/post/search", server.SearchPost)

	// ログイン後
	authRoutes := r.Group("/").Use(authMiddleware(server.tokenMaker))

	// ユーザー
	authRoutes.DELETE("/users/:id", server.DeleteUser)
	authRoutes.GET("/users/:id", server.GetUserData)
	authRoutes.PUT("/users/:id/password", server.ResetPassword)
	authRoutes.PUT("/users/:id/disease-condition", server.UpdateDiseaseAndCondition)
	authRoutes.PUT("/users/:id/email", server.UpdateEmail)
	authRoutes.PUT("/users/:id/privacy", server.UpdateIsPrivacy)
	authRoutes.PUT("/users/:id/name", server.UpdateName)

	// 投稿
	authRoutes.POST("/post/add", server.CreatePost)
	authRoutes.DELETE("/post/del", server.DeletePost)
	authRoutes.PUT("/post/update", server.UpdatePost)

	// タップ
	r.POST("/tag/add", server.CreateTag)
	r.POST("/tag/get", server.FindTag)

	// Bookmark
	authRoutes.POST("/bookmark/add", server.CreateBookmark)
	authRoutes.DELETE("/bookmark/del", server.DeleteBookmark)
	authRoutes.GET("/bookmark/get", server.GetBookmark)

	// コメント
	r.GET("/getcommentlist/:post_id", server.GetCommentsList)
	authRoutes.POST("/createcomment", server.CreateComment)
	r.PUT("/updatecomment", server.UpdateComments)
	r.DELETE("/deletecomment/:comment_id", server.DeleteComments)

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

var (
	ErrInvalidInput     = errors.New("無効な入力です")
	ErrPermissionDenied = errors.New("権限が拒否されました")
	ErrConflict         = errors.New("リソースの競合です")
	ErrUnauthorized     = errors.New("認証に失敗しました")
)

func errorResponse(err error, msg string) gin.H {
	return gin.H{"エラー": err.Error(), "メッセージ": msg}
}

func handleDBError(ctx *gin.Context, err error, msg string) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		ctx.JSON(http.StatusNotFound, errorResponse(err, msg))
	case errors.Is(err, ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, errorResponse(err, msg))
	case errors.Is(err, ErrPermissionDenied):
		ctx.JSON(http.StatusForbidden, errorResponse(err, msg))
	case errors.Is(err, ErrConflict):
		ctx.JSON(http.StatusConflict, errorResponse(err, msg))
	case errors.Is(err, ErrUnauthorized):
		ctx.JSON(http.StatusUnauthorized, errorResponse(err, msg))
	default:
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, msg))
	}
}
