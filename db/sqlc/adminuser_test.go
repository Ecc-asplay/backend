package db

import (
	"context"
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomAdminUser(t *testing.T) CreateAdminUserParams {
	hash, err := util.Hash(util.RandomPassword(20))
	if err != nil {
		log.Fatal(err)
	}

	Data := CreateAdminUserParams{
		AdminID:      util.CreateUUID(),
		Email:        gofakeit.Email(),
		Hashpassword: hash,
		StaffName:    gofakeit.Name(),
		Department:   gofakeit.JobTitle(),
	}

	admin, err := testQueries.CreateAdminUser(context.Background(), Data)
	require.NoError(t, err)
	require.NotEmpty(t, admin)
	require.NotEmpty(t, admin.AdminID)
	require.Equal(t, Data.Email, admin.Email)
	require.Equal(t, Data.Hashpassword, admin.Hashpassword)
	require.Equal(t, Data.StaffName, admin.StaffName)
	require.Equal(t, Data.Department, admin.Department)

	return Data
}

func TestCreateAdminUser(t *testing.T) {
	CreateRandomAdminUser(t)
}

func TestGetAdminLogin(t *testing.T) {
	pw := "abcd12345"
	hash, err := util.Hash(pw)
	if err != nil {
		log.Fatal(err)
	}

	Data := CreateAdminUserParams{
		AdminID:      util.CreateUUID(),
		Email:        "abcde12345@gmail.com",
		Hashpassword: hash,
		StaffName:    gofakeit.Name(),
		Department:   gofakeit.JobTitle(),
	}

	admin, err := testQueries.CreateAdminUser(context.Background(), Data)
	require.NoError(t, err)
	require.NotEmpty(t, admin)
	require.NotEmpty(t, admin.AdminID)
	require.Equal(t, Data.Email, admin.Email)
	require.Equal(t, Data.Hashpassword, admin.Hashpassword)
	require.Equal(t, Data.StaffName, admin.StaffName)
	require.Equal(t, Data.Department, admin.Department)

	Login, err := testQueries.GetAdminLogin(context.Background(), admin.Email)
	require.NoError(t, err)
	require.NotEmpty(t, Login)
	require.Equal(t, admin.AdminID, Login.AdminID)

	isValid, err := util.CheckPassword(pw, Login.Hashpassword)
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestDeleteAdminUser(t *testing.T) {
	data := CreateRandomAdminUser(t)

	err := testQueries.DeleteAdminUser(context.Background(), data.Email)
	require.NoError(t, err)
}
