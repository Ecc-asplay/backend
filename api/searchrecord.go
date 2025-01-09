package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Ecc-asplay/backend/token"
)

func (s *Server) GetSearchedRecordList(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "admin" {
		handleDBError(ctx, errors.New("401"), "管理者権限がございません")
		return
	}

	records, err := s.store.GetSearchedRecordList(ctx)
	if err != nil {
		handleDBError(ctx, err, "検索レコードの取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, records)
}
