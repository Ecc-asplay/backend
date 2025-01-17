package api

import (
	"encoding/json"
	"errors"
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

// 以下は仮の値
// "active": コメントが有効で、ユーザーに表示可能な状態。
// "flagged": コメントが不適切な内容としてフラグされた状態。

type createCommentRequest struct {
	PostID     uuid.UUID `json:"post_id" binding:"required"`
	Comments   string    `json:"comments" binding:"required"`
	IsPublic   bool      `json:"is_public"`
	IsCensored bool      `json:"is_censored"`
}

func (s *Server) CreateComment(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req createCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメント作成：無効な入力データです")
		return
	}

	arg := db.CreateCommentsParams{
		CommentID:  util.CreateUUID(),
		UserID:     authPayload.UserID,
		PostID:     req.PostID,
		Status:     "active",
		IsPublic:   false,
		Comments:   req.Comments,
		IsCensored: req.IsCensored,
	}

	comment, err := s.store.CreateComments(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "コメント作成に失敗しました")
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

func (s *Server) GetPostCommentsList(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		handleDBError(ctx, errors.New("404"), "コメントリスト取得：トークンない")
		return
	}

	postIDStr := ctx.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		handleDBError(ctx, err, "コメントリスト取得：コメントID取得に失敗しました")
		return
	}

	comments, err := s.store.GetCommentsList(ctx, postID)
	if err != nil {
		handleDBError(ctx, err, "コメントリスト取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

type UpdateCommentRequest struct {
	CommentID uuid.UUID `json:"comment_id" binding:"required"`
	Comments  string    `json:"comments" binding:"required"`
	IsPublic  bool      `json:"is_public"`
}

func (s *Server) UpdateComments(ctx *gin.Context) {
	var req UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメント更新：無効な入力データです")
		return
	}

	arg := db.UpdateCommentsParams{
		CommentID: req.CommentID,
		Status:    "active",
		IsPublic:  req.IsPublic,
		Comments:  req.Comments,
	}

	comment, err := s.store.UpdateComments(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "コメント更新に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func (s *Server) DeleteComments(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		handleDBError(ctx, errors.New("404"), "コメント削除：トークンない")
		return
	}

	commentIDStr := ctx.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		handleDBError(ctx, err, "コメント削除：コメントID取得に失敗しました")
		return
	}

	err = s.store.DeleteComments(ctx, commentID)
	if err != nil {
		handleDBError(ctx, err, "コメント削除に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "コメントが削除されました"})
}

func (s *Server) GetAllComments(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		handleDBError(ctx, errors.New("404"), "全部コメント取得：トークンない")
		return
	}

	allComment, err := s.store.GetAllComments(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "全部コメント取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, allComment)
}

func (s *Server) GetAllPublicComments(ctx *gin.Context) {
	allComments, err := s.redis.Get("allComments").Result()
	if err != nil && err != redis.Nil {
		handleDBError(ctx, err, "Redisコメント取得：データ締め切りました")
		return
	}

	if allComments != "" {
		var comments []db.Comment
		err := json.Unmarshal([]byte(allComments), &comments)
		if err != nil {
			handleDBError(ctx, err, "Redisコメント取得：データ変更を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, comments)
	} else {
		comment, err := s.store.GetAllPublicComments(ctx)
		if err != nil {
			handleDBError(ctx, err, "Psqlコメント取得を失敗しました")
			return
		}

		commentJSON, err := json.Marshal(comment)
		if err != nil {
			handleDBError(ctx, err, "Psqlコメント取得：JSON変更を失敗しました")
			return
		}

		err = s.redis.Set("allComments", commentJSON, 5*time.Minute).Err()
		if err != nil {
			handleDBError(ctx, err, "Redisコメント保存を失敗しました")
			return
		}

		log.Info().Msg("すべてのコメントがRedisに追加されました")
		ctx.JSON(http.StatusOK, comment)
	}
}
