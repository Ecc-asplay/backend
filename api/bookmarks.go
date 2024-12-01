package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/Ecc-asplay/backend/db/sqlc"
)

type createBookmarkRequest struct {
	PostID uuid.UUID `json:"post_id"`
}

func (s *Server) CreateBookmark(ctx *gin.Context) {
	var req createBookmarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	data := db.CreateBookmarksParams{
		PostID: req.PostID,
		// UserID: authPayload.UserID,
	}

	bookmark, err := s.store.CreateBookmarks(ctx, data)
	if err != nil {
		handleDBError(ctx, err)
	}

	ctx.JSON(http.StatusCreated, bookmark)
}
