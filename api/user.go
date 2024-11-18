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
