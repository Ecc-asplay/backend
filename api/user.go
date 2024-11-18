// package api

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	db "github.com/Ecc-asplay/backend/db/sqlc"
// 	"github.com/Ecc-asplay/backend/util"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5/pgtype"
// )

// func (s *Server) CreateUser(ctx *gin.Context) {
// 	var req struct {
// 		Username string      `json:"username" binding:"required"`
// 		Email    string      `json:"email" binding:"required,email"`
// 		Birth    pgtype.Date `json:"birth" binding:"required"`
// 		Gender   string      `json:"gender" binding:"required"`
// 		Password string      `json:"password" binding:"required"`
// 	}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// パスワードをハッシュ化する
// 	hashedPassword, err := util.Hash(req.Password)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	arg := db.CreateUserParams{
// 		UserID:       uuid.New(),
// 		Username:     req.Username,
// 		Email:        req.Email,
// 		Birth:        req.Birth,
// 		Gender:       req.Gender,
// 		Hashpassword: hashedPassword,
// 		Disease:      "",
// 		Condition:    "",
// 	}

// 	user, err := s.store.CreateUser(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, user)
// }

// func (s *Server) DeleteUser(ctx *gin.Context) {
// 	userID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := db.DeleteUserParams{
// 		UserID: userID,
// 		Email:  ctx.Query("email"),
// 	}

// 	err = s.store.DeleteUser(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "user deleted"})
// }

// func (s *Server) GetUserData(ctx *gin.Context) {
// 	userID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := s.store.GetUserData(ctx, userID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, user)
// }

// func (s *Server) ResetPassword(ctx *gin.Context) {
// 	userID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	var req struct {
// 		NewPassword string `json:"new_password" binding:"required"`
// 	}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	hashedPassword, err := util.Hash(req.NewPassword)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	arg := db.ResetPasswordParams{
// 		UserID:          userID,
// 		Hashpassword:    hashedPassword,
// 		ResetPasswordAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
// 	}

// 	err = s.store.ResetPassword(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"status": "password reset successful"})
// }

// func (s *Server) UpdateDiseaseAndCondition(ctx *gin.Context) {}
// func (s *Server) UpdateEmail(ctx *gin.Context)               {}
// func (s *Server) UpdateIsPrivacy(ctx *gin.Context)           {}
// func (s *Server) UpdateName(ctx *gin.Context)                {}
// func (s *Server) LoginUser(ctx *gin.Context)                 {}
package api

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

type Test struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (server *Server) Createuser(ctx *gin.Context) {

}

func (server *Server) CreateUser2(ctx *gin.Context) {
	var req Test

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hash, err := util.Hash(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data := db.CreateUserParams{
		UserID:   util.CreateUUID(),
		Username: gofakeit.Name(),
		Email:    req.Email,
		Birth: pgtype.Date{
			Time:  gofakeit.Date(),
			Valid: true,
		},
		Gender:       util.RandomGender(),
		Disease:      util.RandomDisease(),
		Condition:    util.RandomCondition(),
		Hashpassword: hash,
	}

	user, err := server.store.CreateUser(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
