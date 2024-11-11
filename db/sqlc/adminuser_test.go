package db

import (
	"context"
	"log"
	"testing"

	"github.com/Ecc-asplay/backend/util"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func CreateRandomAdminUser() CreateAdminUserParams {
	hash, err := util.Hash(util.RandomPassword(20))
	if err != nil {
		log.Fatal(err)
	}

	Data := CreateAdminUserParams{
		Email:        gofakeit.Email(),
		Hashpassword: hash,
		StaffName:    gofakeit.Name(),
		Department:   gofakeit.JobTitle(),
	}
	return Data
}

func TestCreateAdminUser(t *testing.T) {
	data := CreateRandomAdminUser()

	admin, err := testQueries.CreateAdminUser(context.Background(), data)
	require.NoError(t, err)
	require.NotEmpty(t, admin)
	require.Equal(t, data.Email, admin.Email)
	require.Equal(t, data.Hashpassword, admin.Hashpassword)
	require.Equal(t, data.StaffName, admin.StaffName)
	require.Equal(t, data.Department, admin.Department)
}

func TestGetPasswordToAdminLogin(t *testing.T) {
	data := CreateRandomAdminUser()
	_, err := testQueries.CreateAdminUser(context.Background(), data)
	if err != nil {
		log.Fatal(err)
	}

	Login, err := testQueries.GetPasswordToAdminLogin(context.Background(), data.Email)
	require.NoError(t, err)
	require.Equal(t, data.Hashpassword, Login)
}

func TestDeleteAdminUser(t *testing.T) {
	data := CreateRandomAdminUser()
	admin, err := testQueries.CreateAdminUser(context.Background(), data)
	if err != nil {
		log.Fatal(err)
	}

	err = testQueries.DeleteAdminUser(context.Background(), admin.Email)
	require.NoError(t, err)
}
