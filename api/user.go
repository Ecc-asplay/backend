package api

import (
	"net/http"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (s *Server) Createuser(ctx *gin.Context) {
	// 確認用記述(後で消す)
	var req db.CreateUserParams
	// ↑使わないほうが良い(ハッシュパスワードなどまだない情報があり、エラーがおこるかも？)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.CreateUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
