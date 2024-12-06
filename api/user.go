package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
)

type CreateUserRequset struct {
	Username string      `json:"username" binding:"required"`
	Email    string      `json:"email" binding:"required"`
	Birth    pgtype.Date `json:"birth" binding:"required"`
	Gender   string      `json:"gender"`
	Password string      `json:"password" binding:"required"`
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequset
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "ユーザー作成：無効な入力データです")
		return
	}

	hashedPassword, err := util.Hash(req.Password)
	if err != nil {
		handleDBError(ctx, err, "ユーザー作成：ハッシュ化を失敗しました")
		return
	}

	data := db.CreateUserParams{
		UserID:       util.CreateUUID(),
		Username:     req.Username,
		Email:        req.Email,
		Birth:        req.Birth,
		Gender:       req.Gender,
		Disease:      "",
		Condition:    "",
		Hashpassword: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, data)
	if err != nil {
		handleDBError(ctx, err, "ユーザー作成を失敗しました")
		return
	}

	accessToken, payload, err := s.tokenMaker.CreateToken(user.UserID, "user", s.config.AccessTokenDuration)
	if err != nil {
		handleDBError(ctx, err, "ユーザー作成：トークン作成を失敗しました")
		return
	}

	tokenData := db.CreateTokenParams{
		ID:          util.CreateUUID(),
		UserID:      payload.UserID,
		AccessToken: accessToken,
		Roles:       payload.Role,
		Status:      "SignUp",
		ExpiresAt:   pgtype.Timestamp{Time: payload.ExpiredAt, Valid: true},
	}

	Token, err := s.store.CreateToken(ctx, tokenData)
	if err != nil {
		handleDBError(ctx, err, "トークン保存を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Access_Token":     accessToken,
		"SignUp_At":        Token.TakeAt,
		"User_Information": user,
	})
}

func (s *Server) DeleteUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.store.GetUserData(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "ユーザー削除：無効な入力データです")
		return
	}

	data := db.DeleteUserParams{
		UserID: authPayload.UserID,
		Email:  user.Email,
	}

	err = s.store.DeleteUser(ctx, data)
	if err != nil {
		handleDBError(ctx, err, "ユーザー削除を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "ユーザーが削除されました"})
}

func (s *Server) GetUserData(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.store.GetUserData(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err, "ユーザー取得を失敗しました")
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (s *Server) ResetPassword(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "パスワード再設定：無効な入力データです")
		return
	}

	hashedPassword, err := util.Hash(req.NewPassword)
	if err != nil {
		handleDBError(ctx, err, "パスワード再設定：ハッシュ化を失敗しました")
		return
	}

	arg := db.ResetPasswordParams{
		UserID:          authPayload.UserID,
		Hashpassword:    hashedPassword,
		ResetPasswordAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	err = s.store.ResetPassword(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "パスワード再設定を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "パスワードのリセットが成功しました"})
}

func (s *Server) UpdateDiseaseAndCondition(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		Disease   string `json:"disease" binding:"required"`
		Condition string `json:"condition" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "病歴と病状更新：無効な入力データです")
		return
	}

	arg := db.UpdateDiseaseAndConditionParams{
		UserID:    authPayload.UserID,
		Disease:   req.Disease,
		Condition: req.Condition,
	}

	err := s.store.UpdateDiseaseAndCondition(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "病歴と病状更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "病歴と病状が正常に更新されました"})
}

func (s *Server) UpdateEmail(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "メール更新：無効な入力データです")
		return
	}

	arg := db.UpdateEmailParams{
		UserID: authPayload.UserID,
		Email:  req.NewEmail,
	}

	err := s.store.UpdateEmail(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "メール更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "メールアドレスが正常に更新されました"})
}

func (s *Server) UpdateIsPrivacy(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		IsPrivacy bool `json:"is_privacy" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "プライバシー更新：無効な入力データです")
		return
	}

	arg := db.UpdateIsPrivacyParams{
		UserID:    authPayload.UserID,
		IsPrivacy: req.IsPrivacy,
	}

	err := s.store.UpdateIsPrivacy(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "プライバシー更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "プライバシー設定が正常に更新されました"})
}

func (s *Server) UpdateName(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req struct {
		NewUsername string `json:"new_username" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "ユーザー名更新：無効な入力データです")
		return
	}

	arg := db.UpdateNameParams{
		UserID:   authPayload.UserID,
		Username: req.NewUsername,
	}

	_, err := s.store.UpdateName(ctx, arg)
	if err != nil {
		handleDBError(ctx, err, "ユーザー名更新を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"状態": "ユーザー名が正常に更新されました"})
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "ログイン：無効な入力データです")
		return
	}

	// ハッシュパスワードを取得
	hashedPassword, err := s.store.GetLogin(ctx, req.Email)
	if err != nil {
		handleDBError(ctx, err, "ログイン：ハッシュパスワード取得を失敗しました")
		return
	}

	// パスワードを検証
	isValid, err := util.CheckPassword(req.Password, hashedPassword.Hashpassword)
	if err != nil {
		handleDBError(ctx, err, "ログイン：パスワード認証を失敗しました")
		return
	}
	if !isValid {
		handleDBError(ctx, err, "無効なメールアドレスまたはパスワード")
		return
	}

	// ユーザー情報を取得
	user, err := s.store.GetUserData(ctx, hashedPassword.UserID)
	if err != nil {
		handleDBError(ctx, err, "ログイン：ユーザーデータ取得を失敗しました")
		return
	}

	accessToken, payload, err := s.tokenMaker.CreateToken(user.UserID, "user", s.config.AccessTokenDuration)
	if err != nil {
		handleDBError(ctx, err, "ログイン：トークン作成を失敗しました")
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
		handleDBError(ctx, err, "ログイン：トークン保存を失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"login_at":     Token.TakeAt,
	})
}
