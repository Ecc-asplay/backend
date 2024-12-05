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

	ctx.JSON(http.StatusCreated, tag)
}

type GetTagRequest struct {
	TagComments string `json:"tag_comments" binding:"required"`
}

func (s *Server) FindTag(ctx *gin.Context) {
	var req GetTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	members, err := s.redis.SMembers("Tag").Result()
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	var result []string
	if len(members) > 0 {
		for _, member := range members {
			if strings.Contains(member, req.TagComments) {
				result = append(result, member)
			}
		}
	}
	if len(result) > 0 {
		ctx.JSON(http.StatusOK, result)
	} else {
		// Psql
		tag, err := s.store.FindTag(ctx, req.TagComments)
		if err != nil {
			handleDBError(ctx, err)
			return
		}

		// Add to key
		err = s.redis.SAdd("Tag", tag).Err()
		if err != nil {
			handleDBError(ctx, err)
			return
		}
		log.Info().Msg("tag added in Redis")

		// TTL
		err = s.redis.Expire("Tag", 1*time.Hour).Err()
		if err != nil {
			handleDBError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, tag)
	}

}
