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
			UserID:         authPayload.UserID,
			CommentID:      req.CommentID,
			CReactionHeart: true,
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
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreateCommentsReactionParams{
			UserID:           authPayload.UserID,
			CommentID:        req.CommentID,
			CReactionHelpful: true,
		}
		reaction, err := s.store.CreateCommentsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "コメント生成Helpful：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
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
	var req UpdateCommentsReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "コメント Reaction Get：無効な入力データです")
		return
	}

	thanks, err := s.store.GetCommentsThanksOfTrue(ctx, req.CommentID)
	if err != nil {
		handleDBError(ctx, err, "コメント Reaction Get thanks：取得を失敗しました")
		return
	}
	heart, err := s.store.GetCommentsHeartOfTrue(ctx, req.CommentID)
	if err != nil {
		handleDBError(ctx, err, "コメント Reaction Get heart：取得を失敗しました")
		return
	}
	helpful, err := s.store.GetCommentsHelpfulOfTrue(ctx, req.CommentID)
	if err != nil {
		handleDBError(ctx, err, "コメント Reaction Get helpful：取得を失敗しました")
		return
	}
	useful, err := s.store.GetCommentsUsefulOfTrue(ctx, req.CommentID)
	if err != nil {
		handleDBError(ctx, err, "コメント Reaction Get Useful：取得を失敗しました")
		return
	}

	reaction := CommentReactionTotals{
		CommentID: req.CommentID,
		Thanks:    thanks,
		Heart:     heart,
		Useful:    useful,
		Helpful:   helpful,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"CommentID": req.CommentID,
		"Reaction":  reaction,
	})
}

func (s *Server) GetAllCommentsReaction(ctx *gin.Context) {
	allCommenReaction, err := s.redis.Get("AllcommentsReacrion").Result()
	if err != nil && err != redis.Nil {
		handleDBError(ctx, err, "RedisコメントReaction取得：データ締め切りました")
		return
	}
	if allCommenReaction != "" {
		var commentsReaction []db.CommentsReaction
		err := json.Unmarshal([]byte(allCommenReaction), &commentsReaction)
		if err != nil {
			handleDBError(ctx, err, "RedisコメントReaction取得：データ変更を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, commentsReaction)

	} else {
		allComment, err := s.redis.Get("allComments").Result()
		if err != nil && err != redis.Nil {
			handleDBError(ctx, err, "Redisコメント取得：データ締め切りました")
			return
		}

		if allComment != "" {
			var comments []db.Comment
			err := json.Unmarshal([]byte(allComment), &comments)
			if err != nil {
				handleDBError(ctx, err, "Redisコメント取得：データ変更を失敗しました")
				return
			}

			var allCommentsReaction []CommentReactionTotals
			for _, comment := range comments {
				thanks, err := s.store.GetCommentsThanksOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get thanks：取得を失敗しました")
					return
				}
				heart, err := s.store.GetCommentsHeartOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get heart：取得を失敗しました")
					return
				}
				helpful, err := s.store.GetCommentsHelpfulOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get helpful：取得を失敗しました")
					return
				}
				useful, err := s.store.GetCommentsUsefulOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get Useful：取得を失敗しました")
					return
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

			commentReactionJSON, err := json.Marshal(allCommentsReaction)
			if err != nil {
				handleDBError(ctx, err, "Psqlコメント取得：JSON変更を失敗しました")
				return
			}

			err = s.redis.Set("AllcommentsReacrion", commentReactionJSON, 5*time.Minute).Err()
			if err != nil {
				handleDBError(ctx, err, "Redis コメント Reaction保存を失敗しました")
				return
			}

			ctx.JSON(http.StatusOK, allCommentsReaction)

		} else {
			comment, err := s.store.GetAllPublicComments(ctx)
			if err != nil {
				handleDBError(ctx, err, "Psqlコメント取得を失敗しました")
				return
			}

			var allCommentsReaction []CommentReactionTotals

			for _, comment := range comment {
				thanks, err := s.store.GetCommentsThanksOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get thanks：取得を失敗しました")
					return
				}
				heart, err := s.store.GetCommentsHeartOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get heart：取得を失敗しました")
					return
				}
				helpful, err := s.store.GetCommentsHelpfulOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get helpful：取得を失敗しました")
					return
				}
				useful, err := s.store.GetCommentsUsefulOfTrue(ctx, comment.CommentID)
				if err != nil {
					handleDBError(ctx, err, "All コメント Reaction Get Useful：取得を失敗しました")
					return
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

			commentReactionJSON, err := json.Marshal(allCommentsReaction)
			if err != nil {
				handleDBError(ctx, err, "Psqlコメント取得：JSON変更を失敗しました")
				return
			}

			err = s.redis.Set("AllcommentsReacrion", commentReactionJSON, 5*time.Minute).Err()
			if err != nil {
				handleDBError(ctx, err, "Redis コメント Reaction保存を失敗しました")
				return
			}

			ctx.JSON(http.StatusOK, allCommentsReaction)
		}
	}
}
