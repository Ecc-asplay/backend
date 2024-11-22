package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Ecc-asplay/backend/util"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	pgDate := pgtype.Date{
		Time:  util.RandomDate(),
		Valid: true,
	}

	hash, err := util.Hash(util.RandomPassword(20))
	if err != nil {
		log.Fatal(err)
	}

	userData := CreateUserParams{
		UserID:       util.CreateUUID(),
		Username:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		Birth:        pgDate,
		Gender:       util.RandomGender(),
		Disease:      util.RandomDisease(),
		Condition:    util.RandomCondition(),
		Hashpassword: hash,
	}

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
	require.NotEmpty(t, user.CreatedAt.Time)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestDeleteUser(t *testing.T) {
	user := CreateRandomUser(t)

	del := DeleteUserParams{
		UserID: user.UserID,
		Email:  user.Email,
	}

	err := testQueries.DeleteUser(context.Background(), del)
	require.NoError(t, err)
}

func TestGetPasswordToUserLogin(t *testing.T) {
	pgDate := pgtype.Date{
		Time:  util.RandomDate(),
		Valid: true,
	}

	pw := util.RandomPassword(20)

	hash, err := util.Hash(pw)
	if err != nil {
		log.Fatal(err)
	}

	userData := CreateUserParams{
		UserID:       util.CreateUUID(),
		Username:     gofakeit.Name(),
		Email:        gofakeit.Email(),
		Birth:        pgDate,
		Gender:       util.RandomGender(),
		Disease:      util.RandomDisease(),
		Condition:    util.RandomCondition(),
		Hashpassword: hash,
	}

	user, err := testQueries.CreateUser(context.Background(), userData)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	takeHash, err := testQueries.GetLogin(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, takeHash)

	checked, err := util.CheckPassword(pw, takeHash.Hashpassword)
	require.NoError(t, err)
	require.True(t, checked)
}

func TestGetUserData(t *testing.T) {
	user := CreateRandomUser(t)

	data, err := testQueries.GetUserData(context.Background(), user.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	require.Equal(t, data.UserID, user.UserID)
	require.Equal(t, data.Username, user.Username)
	require.Equal(t, data.Email, user.Email)
	require.Equal(t, data.Birth.Time, user.Birth.Time)
	require.Equal(t, data.Gender, user.Gender)
	require.Equal(t, data.Disease, user.Disease)
	require.Equal(t, data.Condition, user.Condition)
	require.Equal(t, data.Hashpassword, user.Hashpassword)
	require.Equal(t, data.CreatedAt, user.CreatedAt)
}

func TestResetPassword(t *testing.T) {
	user := CreateRandomUser(t)

	newHash, err := util.Hash(util.RandomPassword(20))
	if err != nil {
		log.Fatal(err)
	}
	pgTime := pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}

	resetPw := ResetPasswordParams{
		UserID:          user.UserID,
		Hashpassword:    newHash,
		ResetPasswordAt: pgTime,
	}

	err = testQueries.ResetPassword(context.Background(), resetPw)
	require.NoError(t, err)
	require.NotEqual(t, user.Hashpassword, resetPw.Hashpassword)
}

func TestUpdateDiseaseAndCondition(t *testing.T) {
	user := CreateRandomUser(t)

	newDate := UpdateDiseaseAndConditionParams{
		UserID:    user.UserID,
		Disease:   util.RandomDisease(),
		Condition: util.RandomCondition(),
	}

	err := testQueries.UpdateDiseaseAndCondition(context.Background(), newDate)
	require.NoError(t, err)
}

func TestUpdateEmail(t *testing.T) {
	user := CreateRandomUser(t)

	newEmail := UpdateEmailParams{
		UserID: user.UserID,
		Email:  gofakeit.Email(),
	}

	err := testQueries.UpdateEmail(context.Background(), newEmail)
	require.NoError(t, err)
	require.NotEqual(t, user.Email, newEmail.Email)
}

func TestUpdateIsPrivacy(t *testing.T) {
	user := CreateRandomUser(t)

	newPrivacy := UpdateIsPrivacyParams{
		UserID:    user.UserID,
		IsPrivacy: !user.IsPrivacy,
	}

	err := testQueries.UpdateIsPrivacy(context.Background(), newPrivacy)
	require.NoError(t, err)
	require.NotEqual(t, user.IsPrivacy, newPrivacy.IsPrivacy)
}

func TestUpdateName(t *testing.T) {
	user := CreateRandomUser(t)

	newName := UpdateNameParams{
		UserID:   user.UserID,
		Username: gofakeit.Name(),
	}
	newNameData, err := testQueries.UpdateName(context.Background(), newName)
	require.NoError(t, err)
	require.NotEmpty(t, newNameData)
	require.NotEqual(t, user.Username, newNameData.Username)
}
