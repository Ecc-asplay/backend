package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	db "github.com/Ecc-asplay/backend/db/sqlc"
)

type CreatePostRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	PostID   uuid.UUID `json:"post_id"`
	ShowID   string    `json:"show_id"`
	Title    string    `json:"title"`
	Feel     string    `json:"feel"`
	Content  []byte    `json:"content"`
	Reaction int32     `json:"reaction"`
	Status   string    `json:"status"`
}

func (s *Server) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postData := db.CreatePostParams{
		UserID:   req.UserID,
		PostID:   req.PostID,
		ShowID:   req.ShowID,
		Title:    req.Title,
		Feel:     req.Feel,
		Content:  req.Content,
		Reaction: req.Reaction,
		Status:   req.Status,
	}

	post, err := s.store.CreatePost(ctx, postData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// Get
func (s *Server) GetAllPost(ctx *gin.Context) {
	allPosts, err := s.redis.Get("AllPosts").Result()
	if err != nil && err != redis.Nil {
		handleDBError(ctx, err)
		return
	}

	if allPosts != "" {
		var posts []db.Post
		err := json.Unmarshal([]byte(allPosts), &posts)
		if err != nil {
			handleDBError(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, posts)
	} else {
		post, err := s.store.GetPostsList(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		postJSON, err := json.Marshal(post)
		if err != nil {
			handleDBError(ctx, err)
			return
		}

		err = s.redis.Set("AllPosts", postJSON, 30*time.Second).Err()
		if err != nil {
			handleDBError(ctx, err)
			return
		}

		log.Info().Msg("all posts added in Redis")

		ctx.JSON(http.StatusOK, post)
	}
}

func (s *Server) FindPost(ctx *gin.Context) {
	var req string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	findPost, err := s.store.GetPostOfKeywords(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, findPost)
}

// Delete
type DeletePostRequest struct {
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
}

func (s *Server) DeletePost(ctx *gin.Context) {
	var req DeletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeletePost(ctx, db.DeletePostParams{
		UserID: req.UserID,
		PostID: req.PostID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
