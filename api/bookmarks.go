package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
)

type bookmarkRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

func (s *Server) CreateBookmark(ctx *gin.Context) {
	var req bookmarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	data := db.CreateBookmarksParams{
		PostID: req.PostID,
		UserID: authPayload.UserID,
	}

	createBookmark, err := s.store.CreateBookmarks(ctx, data)
	if err != nil {
		handleDBError(ctx, err)
	}

	ctx.JSON(http.StatusCreated, createBookmark)
}

func (s *Server) DeleteBookmark(ctx *gin.Context) {
	var req bookmarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	data := db.DeleteBookmarksParams{
		PostID: req.PostID,
		UserID: authPayload.UserID,
	}

	err := s.store.DeleteBookmarks(ctx, data)
	if err != nil {
		handleDBError(ctx, err)
	}

	ctx.Status(http.StatusOK)
}

func (s *Server) GetBookmark(ctx *gin.Context) {

	var req bookmarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	bookmark, err := s.store.GetAllBookmarks(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err)
	}

	ctx.JSON(http.StatusOK, bookmark)
}
