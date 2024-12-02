package api

import (
	"net/http"

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

func (s *Server) GetTag(ctx *gin.Context) {
	var req GetTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Info().Msg(req.TagComments)

	tag, err := s.store.GetTag(ctx, req.TagComments)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tag)
}
