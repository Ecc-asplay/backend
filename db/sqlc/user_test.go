package db

import (
	"context"
	"log"
	"testing"

	"github.com/Ecc-asplay/backend/util"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser() CreateUserParams {
	pgDate := pgtype.Date{
		Time:  util.RandomDate(),
		Valid: true,
	}

	hash, err := util.Hash(util.RandomPassword(20))
	if err != nil {
		log.Fatal(err)
	}

	userData := CreateUserParams{
		UserID:       util.CreateUUID("user"),
		Username:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		Birth:        pgDate,
		Gender:       util.RandomGender(),
		Disease:      util.RandomDisease(),
		Condition:    util.RandomCondition(),
		Hashpassword: hash,
	}

	return userData
}

func TestCreateUser(t *testing.T) {
	userData := CreateRandomUser()
	user, err := testQueries.CreateUser(context.Background(), userData)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.UserID, userData.UserID)
	require.Equal(t, user.Username, userData.Username)
	require.Equal(t, user.Email, userData.Email)
	require.Equal(t, user.Birth.Time, userData.Birth.Time)
	require.Equal(t, user.Gender, userData.Gender)
	require.Equal(t, user.Disease, userData.Disease)
	require.Equal(t, user.Condition, userData.Condition)
	require.Equal(t, user.Hashpassword, userData.Hashpassword)
}
