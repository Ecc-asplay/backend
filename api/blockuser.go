package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
)

type CreateBlockUserRequest struct {
	BlockUserID uuid.UUID `json:"block_user_id" binding:"required"`
	Reason      string    `json:"reason" binding:"required"`
}

func (s *Server) CreateBlockUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req CreateBlockUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "無効な入力データです")
		return
	}

	arg := db.CreateBlockParams{
		UserID:      authPayload.UserID,
		BlockuserID: req.BlockUserID,
		Reason:      req.Reason,
		Status:      "blocked",
	}

	blockedUser, err := s.store.CreateBlock(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "ユーザーのブロックに失敗しました")
		return
	}

	ctx.JSON(http.StatusCreated, blockedUser)
}

func (s *Server) GetBlockUsersByUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	blockedUsers, err := s.store.GetBlockUserlist(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "ブロックしたユーザーの一覧取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, blockedUsers)
}

func (s *Server) GetAllBlockedUsers(ctx *gin.Context) {
	blockedUsers, err := s.store.GetAllBlockUsersList(ctx)
	if err != nil {
		handleDBError(ctx, err, "すべてのブロックユーザーの取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, blockedUsers)
}

type UnblockUserRequest struct {
	BlockUserID uuid.UUID `json:"block_user_id" binding:"required"`
}

func (s *Server) UnblockUser(ctx *gin.Context) {
	var req UnblockUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "無効な入力データです")
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UnBlockUserParams{
		UserID:      authPayload.UserID,
		BlockuserID: req.BlockUserID,
		Status:      "unblocked",
	}

	unblockedUser, err := s.store.UnBlockUser(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "ユーザーのブロック解除に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, unblockedUser)
}
