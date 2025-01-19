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

type UpdatePostReactionRequest struct {
	PostID uuid.UUID `json:"post_id" binding:"required"`
}

func (s *Server) UpdatePostReactionThanks(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdatePostReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿Thanks：無効な入力データです")
		return
	}

	data := db.UpdatePostsReactionThanksParams{
		PostID: req.PostID,
		UserID: authPayload.UserID,
	}

	thanks, err := s.store.UpdatePostsReactionThanks(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreatePostsReactionParams{
			UserID:          authPayload.UserID,
			PostID:          req.PostID,
			PReactionThanks: true,
		}
		reaction, err := s.store.CreatePostsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "投稿生成Thanks：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "投稿Thanks：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, thanks)
}

func (s *Server) UpdatePostReactionHeart(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdatePostReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿Heart：無効な入力データです")
		return
	}

	data := db.UpdatePostsReactionHeartParams{
		PostID: req.PostID,
		UserID: authPayload.UserID,
	}

	heart, err := s.store.UpdatePostsReactionHeart(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreatePostsReactionParams{
			UserID:         authPayload.UserID,
			PostID:         req.PostID,
			PReactionHeart: true,
		}
		reaction, err := s.store.CreatePostsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "投稿生成Heart：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "投稿Heart：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, heart)
}

func (s *Server) UpdatePostReactionUesful(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdatePostReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿Useful：無効な入力データです")
		return
	}

	data := db.UpdatePostsReactionUsefulParams{
		PostID: req.PostID,
		UserID: authPayload.UserID,
	}

	useful, err := s.store.UpdatePostsReactionUseful(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreatePostsReactionParams{
			UserID:          authPayload.UserID,
			PostID:          req.PostID,
			PReactionUseful: true,
		}
		reaction, err := s.store.CreatePostsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "投稿生成Useful：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "投稿Useful：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, useful)
}

func (s *Server) UpdatePostReactionHelpful(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req UpdatePostReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿Helpful：無効な入力データです")
		return
	}

	data := db.UpdatePostsReactionHelpfulParams{
		PostID: req.PostID,
		UserID: authPayload.UserID,
	}

	useful, err := s.store.UpdatePostsReactionHelpful(ctx, data)
	if errors.Is(err, sql.ErrNoRows) {
		data := db.CreatePostsReactionParams{
			UserID:           authPayload.UserID,
			PostID:           req.PostID,
			PReactionHelpful: true,
		}
		reaction, err := s.store.CreatePostsReaction(ctx, data)
		if err != nil {
			handleDBError(ctx, err, "投稿生成Helpful：登録を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, reaction)
	} else {
		handleDBError(ctx, err, "投稿Helpful：更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, useful)
}

type PostReactionTotals struct {
	PostID  uuid.UUID `json:"post_id"`
	Thanks  int64     `json:"thanks`
	Heart   int64     `json:"heart"`
	Useful  int64     `json:"useful"`
	Helpful int64     `json:"helpful"`
}

func (s *Server) GetPostReactions(ctx *gin.Context) {
	var req UpdatePostReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "投稿Reaction：無効な入力データです")
		return
	}

	thanks, err := s.store.GetPostsThanksOfTrue(ctx, req.PostID)
	if err != nil {
		handleDBError(ctx, err, "投稿Reaction：thanks取得を失敗しました")
		return
	}
	heart, err := s.store.GetPostsHeartOfTrue(ctx, req.PostID)
	if err != nil {
		handleDBError(ctx, err, "投稿Reaction：heart取得を失敗しました")
		return
	}
	helpful, err := s.store.GetPostsHelpfulOfTrue(ctx, req.PostID)
	if err != nil {
		handleDBError(ctx, err, "投稿Reaction：helpful取得を失敗しました")
		return
	}
	useful, err := s.store.GetPostsUsefulOfTrue(ctx, req.PostID)
	if err != nil {
		handleDBError(ctx, err, "投稿Reaction：Useful取得を失敗しました")
		return
	}

	reaction := PostReactionTotals{
		PostID:  req.PostID,
		Thanks:  thanks,
		Heart:   heart,
		Useful:  useful,
		Helpful: helpful,
	}

	ctx.JSON(http.StatusOK, reaction)
}

func (s *Server) GetAllPostsReaction(ctx *gin.Context) {
	// Redisから投稿Reactionデータ取り
	allPostsReacrion, err := s.redis.Get("allPostsReacrion").Result()
	if err != nil && err != redis.Nil {
		handleDBError(ctx, err, "Redis投稿Reaction取得：データ締め切りました")
		return
	}

	// Redisに投稿Reactionデータあり
	if allPostsReacrion != "" {
		var postsReaction []PostReactionTotals
		err := json.Unmarshal([]byte(allPostsReacrion), &postsReaction)
		if err != nil {
			handleDBError(ctx, err, "Redis投稿Reaction取得：データ変更を失敗しました")
			return
		}

		ctx.JSON(http.StatusOK, postsReaction)
	} else {
		// Redisに投稿Reactionデータない
		allPosts, err := s.redis.Get("allPosts").Result()
		if err != nil && err != redis.Nil {
			handleDBError(ctx, err, "Redis投稿取得：データ締め切りました")
			return
		}

		if allPosts != "" {
			// Redisに投稿データあり
			var posts []db.Post
			err := json.Unmarshal([]byte(allPosts), &posts)
			if err != nil {
				handleDBError(ctx, err, "Redis投稿取得：データ変更を失敗しました")
				return
			}

			allPostsReaction := TakeAllPostsReaction(ctx, s, posts)
			PostSaveToRedis(ctx, s, allPostsReaction, "allPostsReacrion")

			ctx.JSON(http.StatusOK, allPostsReaction)
		} else {
			// Redisに投稿データない
			post, err := s.store.GetPostsList(ctx)
			if err != nil {
				handleDBError(ctx, err, "Psql投稿取得を失敗しました")
				return
			}
			PostSaveToRedis(ctx, s, post, "allPosts")

			allPostsReaction := TakeAllPostsReaction(ctx, s, post)
			PostSaveToRedis(ctx, s, allPostsReaction, "allPostsReacrion")

			ctx.JSON(http.StatusOK, allPostsReaction)
		}
	}
}

func TakeAllPostsReaction(ctx *gin.Context, s *Server, post []db.Post) []PostReactionTotals {
	var allPostsReaction []PostReactionTotals
	for _, post := range post {
		thanks, err := s.store.GetPostsThanksOfTrue(ctx, post.PostID)
		if err != nil {
			handleDBError(ctx, err, "投稿Reaction：thanks取得を失敗しました")
			return nil
		}
		heart, err := s.store.GetPostsHeartOfTrue(ctx, post.PostID)
		if err != nil {
			handleDBError(ctx, err, "投稿Reaction：heart取得を失敗しました")
			return nil
		}
		helpful, err := s.store.GetPostsHelpfulOfTrue(ctx, post.PostID)
		if err != nil {
			handleDBError(ctx, err, "投稿Reaction：helpful取得を失敗しました")
			return nil
		}
		useful, err := s.store.GetPostsUsefulOfTrue(ctx, post.PostID)
		if err != nil {
			handleDBError(ctx, err, "投稿Reaction：Useful取得を失敗しました")
			return nil
		}

		reaction := PostReactionTotals{
			PostID:  post.PostID,
			Thanks:  thanks,
			Heart:   heart,
			Useful:  useful,
			Helpful: helpful,
		}

		allPostsReaction = append(allPostsReaction, reaction)
	}

	return allPostsReaction
}

func PostSaveToRedis(ctx *gin.Context, s *Server, data any, tagname string) {
	JSON, err := json.Marshal(data)
	if err != nil {
		handleDBError(ctx, err, "投稿Reaction保存：JSON変更を失敗しました")
	}

	err = s.redis.Set(tagname, JSON, 5*time.Minute).Err()
	if err != nil {
		handleDBError(ctx, err, "Redis投稿Reaction保存を失敗しました")
	}
}
