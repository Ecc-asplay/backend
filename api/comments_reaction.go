package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
)

type UpdateCommentsReactionRequest struct {
	CommentID uuid.UUID `json:"comment_id" binding:"required"`
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data := db.CreateCommentsReactionParams{
				UserID:           authPayload.UserID,
				CommentID:        req.CommentID,
				CReactionThanks:  true,
				CReactionHelpful: false,
				CReactionUseful:  false,
				CReactionHeart:   false,
			}
			reaction, err := s.store.CreateCommentsReaction(ctx, data)
			if err != nil {
				handleDBError(ctx, err, "コメント生成Thanks：登録を失敗しました")
				return
			}

			ctx.JSON(http.StatusOK, reaction)
		}
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data := db.CreateCommentsReactionParams{
				UserID:           authPayload.UserID,
				CommentID:        req.CommentID,
				CReactionHeart:   true,
				CReactionHelpful: false,
				CReactionUseful:  false,
				CReactionThanks:  false,
			}
			reaction, err := s.store.CreateCommentsReaction(ctx, data)
			if err != nil {
				handleDBError(ctx, err, "コメント生成Heart：登録を失敗しました")
				return
			}

			ctx.JSON(http.StatusOK, reaction)
		}
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data := db.CreateCommentsReactionParams{
				UserID:           authPayload.UserID,
				CommentID:        req.CommentID,
				CReactionUseful:  true,
				CReactionThanks:  false,
				CReactionHeart:   false,
				CReactionHelpful: false,
			}
			reaction, err := s.store.CreateCommentsReaction(ctx, data)
			if err != nil {
				handleDBError(ctx, err, "コメント生成Useful：登録を失敗しました")
				return
			}

			ctx.JSON(http.StatusOK, reaction)
		}
		handleDBError(ctx, err, "コメントUseful：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, Heart)
}

func (s *Server) UpdateCommentReactionHelpful(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdateCommentsReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメントHelpful：無効な入力データです")
		return
	}

	data := db.UpdateCommentsReactionHelpfulParams{
		CommentID: req.CommentID,
		UserID:    authPayload.UserID,
	}

	Heart, err := s.store.UpdateCommentsReactionHelpful(ctx, data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data := db.CreateCommentsReactionParams{
				UserID:           authPayload.UserID,
				CommentID:        req.CommentID,
				CReactionHelpful: true,
				CReactionThanks:  false,
				CReactionUseful:  false,
				CReactionHeart:   false,
			}
			reaction, err := s.store.CreateCommentsReaction(ctx, data)
			if err != nil {
				handleDBError(ctx, err, "コメント生成Helpful：登録を失敗しました")
				return
			}

			ctx.JSON(http.StatusOK, reaction)
		}
		handleDBError(ctx, err, "コメントHelpful：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, Heart)
}

type CommentReactionTotals struct {
	CommentID uuid.UUID `json:"comment_id"`
	Thanks    int64     `json:"thanks`
	Heart     int64     `json:"heart"`
	Useful    int64     `json:"useful"`
	Helpful   int64     `json:"helpful"`
}

func (s *Server) GetCommentReactions(ctx *gin.Context) {
	commentIDStr := ctx.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		handleDBError(ctx, err, "コメントリスト取得：コメントID取得に失敗しました")
	}

	thanks, err := s.store.GetCommentsThanksOfTrue(ctx, commentID)
	if err != nil {
		handleDBError(ctx, err, "コメントReaction：thanks取得を失敗しました")
		return
	}
	heart, err := s.store.GetCommentsHeartOfTrue(ctx, commentID)
	if err != nil {
		handleDBError(ctx, err, "コメントReaction：heart取得を失敗しました")
		return
	}
	helpful, err := s.store.GetCommentsHelpfulOfTrue(ctx, commentID)
	if err != nil {
		handleDBError(ctx, err, "コメントReaction：helpful取得を失敗しました")
		return
	}
	useful, err := s.store.GetCommentsUsefulOfTrue(ctx, commentID)
	if err != nil {
		handleDBError(ctx, err, "コメントReaction：Useful取得を失敗しました")
		return
	}

	reaction := CommentReactionTotals{
		CommentID: commentID,
		Thanks:    thanks,
		Heart:     heart,
		Useful:    useful,
		Helpful:   helpful,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"CommentID": commentID,
		"Reaction":  reaction,
	})
}

func (s *Server) GetAllCommentsReaction(ctx *gin.Context) {
	// RedisからコメントReactionデータ取り
	allCommenReaction, err := s.redis.Get("allCommentsReacrion").Result()
	if err != nil && err != redis.Nil {
		handleDBError(ctx, err, "RedisコメントReaction取得：データ締め切りました")
		return
	}

	// RedisにコメントReactionデータあり
	if allCommenReaction != "" {
		var commentsReaction []CommentReactionTotals
		err := json.Unmarshal([]byte(allCommenReaction), &commentsReaction)
		if err != nil {
			handleDBError(ctx, err, "RedisコメントReaction取得：データ変更を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, commentsReaction)
	} else {
		// RedisにコメントReactionデータない
		allComment, err := s.redis.Get("allComments").Result()
		if err != nil && err != redis.Nil {
			handleDBError(ctx, err, "Redisコメント取得：データ締め切りました")
			return
		}

		if allComment != "" {
			// Redisにコメントデータあり
			var comments []db.Comment
			err := json.Unmarshal([]byte(allComment), &comments)
			if err != nil {
				handleDBError(ctx, err, "Redisコメント取得：データ変更を失敗しました")
				return
			}

			allCommentsReaction := TakeAllCommentsReaction(ctx, s, comments)
			CommentSaveToRedis(ctx, s, allCommentsReaction, "allCommentsReacrion")

			ctx.JSON(http.StatusOK, allCommentsReaction)
		} else {
			// Redisにコメントデータない
			comment, err := s.store.GetAllPublicComments(ctx)
			if err != nil {
				handleDBError(ctx, err, "Psqlコメント取得を失敗しました")
				return
			}
			CommentSaveToRedis(ctx, s, comment, "allComments")

			allCommentsReaction := TakeAllCommentsReaction(ctx, s, comment)
			CommentSaveToRedis(ctx, s, allCommentsReaction, "allCommentsReacrion")

			ctx.JSON(http.StatusOK, allCommentsReaction)
		}
	}
}

func TakeAllCommentsReaction(ctx *gin.Context, s *Server, comment []db.Comment) []CommentReactionTotals {
	var allCommentsReaction []CommentReactionTotals
	for _, comment := range comment {
		thanks, err := s.store.GetCommentsThanksOfTrue(ctx, comment.CommentID)
		if err != nil {
			handleDBError(ctx, err, "コメントReaction：thanks取得を失敗しました")
			return nil
		}
		heart, err := s.store.GetCommentsHeartOfTrue(ctx, comment.CommentID)
		if err != nil {
			handleDBError(ctx, err, "コメントReaction：heart取得を失敗しました")
			return nil
		}
		helpful, err := s.store.GetCommentsHelpfulOfTrue(ctx, comment.CommentID)
		if err != nil {
			handleDBError(ctx, err, "コメントReaction：helpful取得を失敗しました")
			return nil
		}
		useful, err := s.store.GetCommentsUsefulOfTrue(ctx, comment.CommentID)
		if err != nil {
			handleDBError(ctx, err, "コメントReaction：Useful取得を失敗しました")
			return nil
		}

		reaction := CommentReactionTotals{
			CommentID: comment.CommentID,
			Thanks:    thanks,
			Heart:     heart,
			Useful:    useful,
			Helpful:   helpful,
		}

		allCommentsReaction = append(allCommentsReaction, reaction)
	}

	return allCommentsReaction
}

func CommentSaveToRedis(ctx *gin.Context, s *Server, data any, tagname string) {
	JSON, err := json.Marshal(data)
	if err != nil {
		handleDBError(ctx, err, "コメントReaction保存：JSON変更を失敗しました")
	}

	err = s.redis.Set(tagname, JSON, 5*time.Minute).Err()
	if err != nil {
		handleDBError(ctx, err, "RedisコメントReaction保存を失敗しました")
	}
}
