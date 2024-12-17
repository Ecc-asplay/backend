package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

func CreateAndLoginManagment(t *testing.T) string {
	server := newTestServer(t)
	hash, err := util.Hash("abcde12345")
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	Email := gofakeit.Email()

	data := db.CreateAdminUserParams{
		AdminID:      util.CreateUUID(),
		Email:        Email,
		Hashpassword: hash,
		StaffName:    gofakeit.Name(),
		Department:   gofakeit.Job().Title,
	}

	management, err := server.store.CreateAdminUser(context.Background(), data)
	require.NoError(t, err)
	require.NotEmpty(t, management)

	loginData := LoginRequest{
		Email:    management.Email,
		Password: "abcde12345",
	}

	adminData, err := server.store.GetAdminLogin(context.Background(), loginData.Email)
	require.NoError(t, err)
	require.NotEmpty(t, adminData)

	isValid, err := util.CheckPassword(loginData.Password, adminData.Hashpassword)
	require.NoError(t, err)
	require.True(t, isValid)

	accessToken, payload, err := server.tokenMaker.CreateToken(adminData.AdminID, "admin", server.config.AccessTokenDuration)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, payload)

	tokenData := db.CreateTokenParams{
		ID:          util.CreateUUID(),
		UserID:      payload.UserID,
		AccessToken: accessToken,
		Roles:       payload.Role,
		Status:      "Login",
		ExpiresAt:   pgtype.Timestamp{Time: payload.ExpiredAt, Valid: true},
	}

	token, err := server.store.CreateToken(context.Background(), tokenData)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return accessToken
}

func RandomCreateAdminUser(t *testing.T) db.Adminuser {
	adminToken := `Bearer ` + CreateAndLoginManagment(t)
	adminData := CreateAdminRequest{
		Email:      gofakeit.Email(),
		Password:   util.RandomString(20),
		StaffName:  gofakeit.Name(),
		Department: gofakeit.Job().Title,
	}

	var createAdmin db.Adminuser

	t.Run("RandomCreateAdminUser", func(t *testing.T) {
		recorder := APITestAfterLogin(t, adminData, http.MethodPost, "/admin/create", adminToken)
		require.Equal(t, http.StatusCreated, recorder.Code)
		require.NotEmpty(t, recorder.Body)

		err := json.Unmarshal(recorder.Body.Bytes(), &createAdmin)
		require.NoError(t, err)
		fmt.Println(" ")
	})
	return createAdmin
}
