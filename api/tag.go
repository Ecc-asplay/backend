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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tagData := db.CreateTagParams{
		PostID:      req.PostID,
		TagComments: req.TagComments,
	}

	tag, err := s.store.CreateTag(ctx, tagData)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	// Redisに追加する
	go func(tag string) {
		err := s.redis.SAdd("Tag", tag).Err()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to add Tag to Redis")
			return
		}

		err = s.redis.Expire("Tag", 1*time.Hour).Err()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to set TTL for Tag in Redis")
			return
		}
		log.Info().Msg("Tag successfully updated in Redis")
	}(tag.TagComments)

	ctx.JSON(http.StatusCreated, tag)
}

type GetTagRequest struct {
	TagComments string `json:"tag_comments" binding:"required"`
}

func (s *Server) FindTag(ctx *gin.Context) {
	var req GetTagRequest
	var result []string

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Redisから取る
	members, err := s.redis.SMembers("Tag").Result()
	if err != nil {
		log.Warn().Err(err).Msg("Redis SMembers failed")
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
		handleDBError(ctx, err)
		return
	}

	// Redis更新
	go func(tag []string) {
		err := s.redis.SAdd("Tag", tag).Err()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to add Tag to Redis")
			return
		}

		err = s.redis.Expire("Tag", 1*time.Hour).Err()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to set TTL for Tag in Redis")
			return
		}
		log.Info().Msg("Tag successfully updated in Redis")
	}(tag)

	ctx.JSON(http.StatusOK, tag)
}
