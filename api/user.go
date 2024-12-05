package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/errors"
	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
)

type User struct {
	Username string      `json:"username" binding:"required"`
	Email    string      `json:"email" binding:"required"`
	Birth    pgtype.Date `json:"birth" binding:"required"`
	Gender   string      `json:"gender"`
	Password string      `json:"password" binding:"required"`
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.HandleDBError(ctx, err)
		return
	}

	hashedPassword, err := util.Hash(req.Password)
	if err != nil {
		handleDBError(ctx, err)
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, payload, err := s.tokenMaker.CreateToken(user.UserID, "user", s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	data := db.DeleteUserParams{
		UserID: authPayload.UserID,
		Email:  user.Email,
	}

	err = s.store.DeleteUser(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "user deleted"})
}

func (s *Server) GetUserData(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.store.GetUserData(ctx, authPayload.UserID)
	if err != nil {
		handleDBError(ctx, err)
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.Hash(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ResetPasswordParams{
		UserID:          authPayload.UserID,
		Hashpassword:    hashedPassword,
		ResetPasswordAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	err = s.store.ResetPassword(ctx, arg)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "password reset successful"})
}

func (s *Server) UpdateDiseaseAndCondition(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		Disease   string `json:"disease" binding:"required"`
		Condition string `json:"condition" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateDiseaseAndConditionParams{
		UserID:    authPayload.UserID,
		Disease:   req.Disease,
		Condition: req.Condition,
	}

	err := s.store.UpdateDiseaseAndCondition(ctx, arg)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "disease and condition updated successfully"})
}

func (s *Server) UpdateEmail(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateEmailParams{
		UserID: authPayload.UserID,
		Email:  req.NewEmail,
	}

	err := s.store.UpdateEmail(ctx, arg)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "email updated successfully"})
}

func (s *Server) UpdateIsPrivacy(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req struct {
		IsPrivacy bool `json:"is_privacy" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateIsPrivacyParams{
		UserID:    authPayload.UserID,
		IsPrivacy: req.IsPrivacy,
	}

	err := s.store.UpdateIsPrivacy(ctx, arg)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "privacy setting updated successfully"})
}

func (s *Server) UpdateName(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req struct {
		NewUsername string `json:"new_username" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateNameParams{
		UserID:   authPayload.UserID,
		Username: req.NewUsername,
	}

	_, err := s.store.UpdateName(ctx, arg)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "username updated successfully"})
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// ハッシュパスワードを取得
	hashedPassword, err := s.store.GetLogin(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid email or password")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// パスワードを検証
	isValid, err := util.CheckPassword(req.Password, hashedPassword.Hashpassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to verify password")))
		return
	}
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid email or password")))
		return
	}

	// ユーザー情報を取得
	user, err := s.store.GetUserData(ctx, hashedPassword.UserID)
	if err != nil {
		handleDBError(ctx, err)
		return
	}

	accessToken, payload, err := s.tokenMaker.CreateToken(user.UserID, "user", s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create token")))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"login_at":     Token.TakeAt,
	})
}
