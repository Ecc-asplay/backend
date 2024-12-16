package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
		recode := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(adminData)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		request, err := http.NewRequest(http.MethodPost, "/admin/create", bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", adminToken)

		server.router.ServeHTTP(recode, request)
		require.NotEmpty(t, recode)

		require.Equal(t, http.StatusCreated, recode.Code)
		require.NotEmpty(t, recode.Body)

		err = json.Unmarshal(recode.Body.Bytes(), &createAdmin)
		require.NoError(t, err)
		fmt.Println(" ")
	})
	return createAdmin
}
