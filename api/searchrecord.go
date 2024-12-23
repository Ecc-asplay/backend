package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/Ecc-asplay/backend/db/sqlc"
)

type CreateSearchedRecordRequest struct {
	SearchContent string `json:"search_content" binding:"required"`
	IsUser        bool   `json:"is_user"`
}

func (s *Server) CreateSearchRecord(ctx *gin.Context) {
	var req CreateSearchedRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "検索記録作成：無効な入力データです")
		return
	}

	arg := db.CreateSearchedRecordParams{
		SearchContent: req.SearchContent,
		IsUser:        req.IsUser,
	}

	record, err := s.store.CreateSearchedRecord(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "検索レコードの作成に失敗しました")
		return
	}

	ctx.JSON(http.StatusCreated, record)
}

func (s *Server) GetSearchedRecordList(ctx *gin.Context) {
	records, err := s.store.GetSearchedRecordList(ctx)
	if err != nil {
		handleDBError(ctx, err, "検索レコードの取得に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, records)
}
