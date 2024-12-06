package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

type LoginAdminRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) LoginAdmin(ctx *gin.Context) {
	var req LoginAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "Admin作成：無効な入力データです")
		return
	}

	adminData, err := s.store.GetAdminLogin(ctx, req.Email)
	if err != nil {
		handleDBError(ctx, err, "Admin ログイン：ハッシュパスワード取得を失敗しました")
		return
	}

	isValid, err := util.CheckPassword(req.Password, adminData.Hashpassword)
	if err != nil {
		handleDBError(ctx, err, "Admin ログイン：パスワード認証を失��しました")
		return
	}
	if !isValid {
		handleDBError(ctx, err, "Admin ログイン：無効なメールアドレスまたはパスワード")
		return
	}

	accessToken, payload, err := s.tokenMaker.CreateToken(adminData.AdminID, "admin", s.config.AccessTokenDuration)
	if err != nil {
		handleDBError(ctx, err, "ユーザー ログイン：トークン作成を失敗しました")
		return
	}
	tokenData := db.CreateTokenParams{
		ID:          util.CreateUUID(),
		UserID:      payload.UserID,
		AccessToken: accessToken,
		Roles:       payload.Role,
		Status:      "Login",
		ExpiresAt:   pgtype.Timestamp{Time: payload.ExpiredAt, Valid: true},
	}
	Token, err := s.store.CreateToken(ctx, tokenData)
	if err != nil {
		handleDBError(ctx, err, "ユーザー ログイン：トークン保存を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"login_at":     Token.TakeAt,
	})
}

type CreateAdminRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	StaffName  string `json:"staff_name" binding:"required"`
	Department string `json:"department" binding:"required"`
}

func (s *Server) CreateAdminUser(ctx *gin.Context) {
	var req CreateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "Admin作成：無効な入力データです")
		return
	}

}
