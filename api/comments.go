package api

import (
	"database/sql"
	"net/http"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ?Status項目は必要か?
// 以下は仮の値
// "active": コメントが有効で、ユーザーに表示可能な状態。
// "flagged": コメントが不適切な内容としてフラグされた状態。

type CreateCommentRequest struct {
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	PostID     uuid.UUID `json:"post_id" binding:"required"`
	Comments   string    `json:"comments" binding:"required"`
	IsPublic   bool      `json:"is_public"`
	Reaction   int32     `json:"reaction"`
	IsCensored bool      `json:"is_censored"`
}

// 動作確認済み
func (s *Server) CreateComment(ctx *gin.Context) {
	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCommentsParams{
		CommentID:  util.CreateUUID(),
		UserID:     req.UserID,
		PostID:     req.PostID,
		Status:     "active",
		IsPublic:   false,
		Comments:   req.Comments,
		Reaction:   req.Reaction,
		IsCensored: req.IsCensored,
	}

	comment, err := s.store.CreateComments(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

// 未完成
func (s *Server) GetCommentsList(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	comments, err := s.store.GetCommentsList(ctx, postID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

type UpdateCommentRequest struct {
	CommentID uuid.UUID `json:"comment_id" binding:"required"`
	Comments  string    `json:"comments" binding:"required"`
	IsPublic  bool      `json:"is_public"`
	Reaction  int32     `json:"reaction"`
}

// 動作確認済み
func (s *Server) UpdateComments(ctx *gin.Context) {
	var req UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCommentsParams{
		CommentID: req.CommentID,
		Status:    "active",
		IsPublic:  req.IsPublic,
		Comments:  req.Comments,
		Reaction:  req.Reaction,
	}

	comment, err := s.store.UpdateComments(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func (s *Server) DeleteComments(ctx *gin.Context) {
	commentIDStr := ctx.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = s.store.DeleteComments(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "comment deleted"})
}
