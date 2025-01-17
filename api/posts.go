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
	"github.com/Ecc-asplay/backend/util"
)

// Create
type CreatePostRequest struct {
	ShowID  string `json:"show_id"`
	Title   string `json:"title"`
	Feel    string `json:"feel"`
	Content []byte `json:"content"`
	Status  string `json:"status" binding:"required"`
}

func (s *Server) CreatePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿作成：無効な入力データです")
		return
	}

	postID := util.CreateUUID()
	var showID string

	if req.ShowID == "" {
		showID = postID.String()
	} else {
		showID = req.ShowID
	}

	postData := db.CreatePostParams{
		UserID:  authPayload.UserID,
		PostID:  postID,
		ShowID:  showID,
		Title:   req.Title,
		Feel:    req.Feel,
		Content: req.Content,
		Status:  req.Status,
	}

	post, err := s.store.CreatePost(ctx, postData)
	if err != nil {
		handleDBError(ctx, err, "投稿作成：保存を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// Get Post of User
func (s *Server) GetUserPost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	userPost, err := s.store.GetUserAllPosts(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "ユーザー投稿取得を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, userPost)
}

// Get all
func (s *Server) GetAllPost(ctx *gin.Context) {
	allPosts, err := s.redis.Get("AllPosts").Result()
	if err != nil && err != redis.Nil {
		handleDBError(ctx, err, "Redis投稿取得：データ締め切りました")
		return
	}

	if allPosts != "" {
		var posts []db.Post
		err := json.Unmarshal([]byte(allPosts), &posts)
		if err != nil {
			handleDBError(ctx, err, "Redis投稿取得：データ変更を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, posts)
	} else {
		post, err := s.store.GetPostsList(ctx)
		if err != nil {
			handleDBError(ctx, err, "Psql投稿取得を失敗しました")
			return
		}

		postJSON, err := json.Marshal(post)
		if err != nil {
			handleDBError(ctx, err, "Psql投稿取得：JSON変更を失敗しました")
			return
		}

		err = s.redis.Set("AllPosts", postJSON, 5*time.Minute).Err()
		if err != nil {
			handleDBError(ctx, err, "Redis投稿保存を失敗しました")
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
		handleDBError(ctx, err, "投稿検索：無効な入力データです")
		return
	}

	findPost, err := s.store.SearchPost(ctx, req)
	if err != nil {
		handleDBError(ctx, err, "投稿検索を失敗しました")
		return
	} else {
		arg := db.CreateSearchedRecordParams{
			SearchContent: req,
			IsUser:        false,
		}

		_, err := s.store.CreateSearchedRecord(ctx, arg)
		if err != nil {
			handleDBError(ctx, err, "投稿検索：検索レコードの作成に失敗しました")
			return
		}
	}

	ctx.JSON(http.StatusOK, findPost)
}

// Delete
type DeletePostRequest struct {
	PostID uuid.UUID `json:"post_id" binding:"required"`
}

func (s *Server) DeletePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req DeletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿削除：無効な入力データです")
		return
	}

	err := s.store.DeletePost(ctx, db.DeletePostParams{
		UserID: authPayload.UserID,
		PostID: req.PostID,
	})

	if err != nil {
		handleDBError(ctx, err, "投稿削除を失敗しました")
		return
	}

	ctx.Status(http.StatusOK)
}

// Update
type UpdatePostsRequest struct {
	PostID      uuid.UUID `json:"post_id" binding:"required"`
	ShowID      string    `json:"show_id"`
	Title       string    `json:"title"`
	Feel        string    `json:"feel"`
	Content     []byte    `json:"content"`
	IsSensitive bool      `json:"is_sensitive"`
	Status      string    `json:"status" binding:"required"`
}

func (s *Server) UpdatePost(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req UpdatePostsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿更新：無効な入力データです")
		return
	}

	newPostData := db.UpdatePostsParams{
		UserID:      authPayload.UserID,
		PostID:      req.PostID,
		ShowID:      req.ShowID,
		Title:       req.Title,
		Feel:        req.Feel,
		Content:     req.Content,
		IsSensitive: req.IsSensitive,
		Status:      req.Status,
	}

	newPost, err := s.store.UpdatePosts(ctx, newPostData)
	if err != nil {
		handleDBError(ctx, err, "投稿更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, newPost)
}
