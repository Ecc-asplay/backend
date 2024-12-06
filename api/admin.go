package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
)

func (s *Server) LoginAdmin(ctx *gin.Context) {
	var req LoginRequest
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "admin" {
		handleDBError(ctx, errors.New("401"), "管理者権限がございません")
		return
	}

	var req CreateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "Admin作成：無効な入力データです")
		return
	}

	hash, err := util.Hash(req.Password)
	if err != nil {
		handleDBError(ctx, err, "Admin作成：ハッシュ化を失敗しました")
		return
	}

	data := db.CreateAdminUserParams{
		AdminID:      util.CreateUUID(),
		Email:        req.Email,
		Hashpassword: hash,
		StaffName:    req.StaffName,
		Department:   req.Department,
	}

	admin, err := s.store.CreateAdminUser(ctx, data)
	if err != nil {
		handleDBError(ctx, err, "Admin作成：登録に失敗しました")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"admin":    admin,
		"password": req.Password,
	})
}

type DeleteAdminUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) DeleteAdminUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != "admin" {
		handleDBError(ctx, errors.New("401"), "管理者権限がございません")
		return
	}

	var req DeleteAdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "Admin削除：無効な入力データです")
		return
	}

	err := s.store.DeleteAdminUser(ctx, req.Email)
	if err != nil {
		handleDBError(ctx, err, "Admin削除を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "Adminが削除されました"})
}
