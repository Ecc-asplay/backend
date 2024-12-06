package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	db "github.com/Ecc-asplay/backend/db/sqlc"
)

type CreateTagRequest struct {
	PostID      uuid.UUID `json:"post_id"`
	TagComments string    `json:"tag_comments" binding:"required"`
}

func (s *Server) CreateTag(ctx *gin.Context) {
	var req CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "タップ作成：無効な入力データです")
		return
	}

	tagData := db.CreateTagParams{
		PostID:      req.PostID,
		TagComments: req.TagComments,
	}

	tag, err := s.store.CreateTag(ctx, tagData)
	if err != nil {
		handleDBError(ctx, err, "Psqlタップ作成を失敗しました")
		return
	}

	// Redisに追加する
	err = s.redis.SAdd("Tag", tag).Err()
	if err != nil {
		handleDBError(ctx, err, "Redisタップ作成を失敗しました")
		return
	}

	err = s.redis.Expire("Tag", 1*time.Hour).Err()
	if err != nil {
		handleDBError(ctx, err, "Redisタップ作成：有効時間設定を失敗しました")
		return
	}

	log.Info().Msg("タグがRedisで正常に更新されました")
	ctx.JSON(http.StatusCreated, tag)
}

type GetTagRequest struct {
	TagComments string `json:"tag_comments" binding:"required"`
}

func (s *Server) FindTag(ctx *gin.Context) {
	var req GetTagRequest
	var result []string

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "タップ検索：無効な入力データです")
		return
	}

	// Redisから取る
	members, err := s.redis.SMembers("Tag").Result()
	if err != nil {
		handleDBError(ctx, err, "Redisタップ検索：データ取得を失敗しました")
		members = nil
	}
	if len(members) > 0 {
		for _, member := range members {
			if strings.Contains(member, req.TagComments) {
				result = append(result, member)
			}
		}
		if len(result) > 0 {
			ctx.JSON(http.StatusOK, result)
			return
		}
	}

	// Psql から取る
	tag, err := s.store.FindTag(ctx, req.TagComments)
	if err != nil {
		handleDBError(ctx, err, "Psqlタップ検索を失敗しました")
		return
	}

	// Redis更新
	err = s.redis.SAdd("Tag", tag).Err()
	if err != nil {
		handleDBError(ctx, err, "Redisタップ検索の保存を失敗しました")
		return
	}

	err = s.redis.Expire("Tag", 1*time.Hour).Err()
	if err != nil {
		handleDBError(ctx, err, "Redisタップ検索：有効時間設定を失敗しました")
		return
	}

	log.Warn().Err(err).Msg("タグのTTLをRedisに設定できました")
	ctx.JSON(http.StatusOK, tag)
}
