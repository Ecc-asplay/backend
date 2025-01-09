package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
)

type CreateNotificationRequest struct {
	Content string `json:"content" binding:"required"`
	Icon    []byte `json:"icon"`
}

func (s *Server) CreateNotification(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "通知作成：無効な入力データです")
		return
	}

	arg := db.CreateNotificationParams{
		UserID:  authPayload.UserID,
		Content: req.Content,
		Icon:    req.Icon,
	}

	notification, err := s.store.CreateNotification(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "通知の作成に失敗しました")
		return
	}

	ctx.JSON(http.StatusCreated, notification)
}

func (s *Server) GetNotificationsByUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	notifications, err := s.store.GetNotification(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "通知の取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func (s *Server) MarkNotificationsAsRead(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	notifications, err := s.store.UpdateNotification(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "通知を既読にできませんでした")
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}
