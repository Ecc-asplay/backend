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
	"github.com/Ecc-asplay/backend/token"
)

// Create
type CreatePostRequest struct {
	PostID   uuid.UUID `json:"post_id"`
	ShowID   string    `json:"show_id"`
	Title    string    `json:"title"`
	Feel     string    `json:"feel"`
	Content  []byte    `json:"content"`
	Reaction int32     `json:"reaction"`
	Status   string    `json:"status"`
}

func (s *Server) CreatePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err)
		return
	}

	postData := db.CreatePostParams{
		UserID:   authPayload.UserID,
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

// Get Post of User
func (s *Server) GetPost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	userPost, err := s.store.GetUserAllPosts(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, userPost)
}

// Get all
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

		err = s.redis.Set("AllPosts", postJSON, 5*time.Minute).Err()
		if err != nil {
			handleDBError(ctx, err)
			return
		}

		log.Info().Msg("すべての投稿がRedisに追加されました")
		ctx.JSON(http.StatusOK, post)
	}
}

// Search
func (s *Server) SearchPost(ctx *gin.Context) {
	var req string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err)
		return
	}

	findPost, err := s.store.SearchPost(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, findPost)
}

// Delete
type DeletePostRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

func (s *Server) DeletePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req DeletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err)
		return
	}

	err := s.store.DeletePost(ctx, db.DeletePostParams{
		UserID: authPayload.UserID,
		PostID: req.PostID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

// Update
type UpdatePostsRequest struct {
	PostID      uuid.UUID `json:"post_id"`
	ShowID      string    `json:"show_id"`
	Title       string    `json:"title"`
	Feel        string    `json:"feel"`
	Content     []byte    `json:"content"`
	Reaction    int32     `json:"reaction"`
	IsSensitive bool      `json:"is_sensitive"`
}

func (s *Server) UpdatePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req UpdatePostsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err)
		return
	}
	newPostData := db.UpdatePostsParams{
		UserID:      authPayload.UserID,
		PostID:      req.PostID,
		ShowID:      req.ShowID,
		Title:       req.Title,
		Feel:        req.Feel,
		Content:     req.Content,
		Reaction:    req.Reaction,
		IsSensitive: req.IsSensitive,
	}

	newPost, err := s.store.UpdatePosts(ctx, newPostData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newPost)
}
