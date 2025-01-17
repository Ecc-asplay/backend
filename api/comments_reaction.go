package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
)

type UpdateCommentsReactionRequest struct {
	CommentID uuid.UUID `json:"commet_id" binding:"required"`
}

func (s *Server) UpdateCommentReactionThanks(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdateCommentsReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメントThanks：無効な入力データです")
		return
	}

	data := db.UpdateCommentsReactionThanksParams{
		CommentID: req.CommentID,
		UserID:    authPayload.UserID,
	}

	thanks, err := s.store.UpdateCommentsReactionThanks(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreateCommentsReactionParams{
			UserID:          authPayload.UserID,
			CommentID:       req.CommentID,
			CReactionThanks: true,
		}
		reaction, err := s.store.CreateCommentsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "コメント生成Thanks：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "コメントThanks：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, thanks)
}

func (s *Server) UpdateCommentReactionHeart(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdateCommentsReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメントHeart：無効な入力データです")
		return
	}

	data := db.UpdateCommentsReactionHeartParams{
		CommentID: req.CommentID,
		UserID:    authPayload.UserID,
	}

	Heart, err := s.store.UpdateCommentsReactionHeart(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreateCommentsReactionParams{
			UserID:          authPayload.UserID,
			CommentID:       req.CommentID,
			CReactionThanks: true,
		}
		reaction, err := s.store.CreateCommentsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "コメント生成Heart：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "コメントHeart：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, Heart)
}

func (s *Server) UpdateCommentReactionUesful(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdateCommentsReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメントUseful：無効な入力データです")
		return
	}

	data := db.UpdateCommentsReactionUsefulParams{
		CommentID: req.CommentID,
		UserID:    authPayload.UserID,
	}

	Heart, err := s.store.UpdateCommentsReactionUseful(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreateCommentsReactionParams{
			UserID:          authPayload.UserID,
			CommentID:       req.CommentID,
			CReactionUseful: true,
		}
		reaction, err := s.store.CreateCommentsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "コメント生成Useful：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "コメントUseful：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, Heart)
}

func (s *Server) UpdateCommentReactionHelpful(ctx *gin.Context) {
}
