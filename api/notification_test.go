package api

import (
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

func RandomNotificationAPI(t *testing.T, user UserRsp) {
	data := db.CreateNotificationParams{
		UserID:  user.User_Information.UserID,
		Content: util.RandomString(20),
		Icon:    gofakeit.ImageJpeg(10, 10),
	}

	server := newTestServer(t)
	require.NotEmpty(t, server)

	_, err := server.store.CreateNotification(context.Background(), data)
	require.NoError(t, err)
}

func TestGetNotificationAPI(t *testing.T) {
	user := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
	})

	token := "Bearer " + user.Access_Token
	for i := 0; i < 20; i++ {
		RandomNotificationAPI(t, user)
	}

	testCases := []struct {
		name          string
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var notifications []db.Notification
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&notifications))
				require.GreaterOrEqual(t, len(notifications), 1)
			},
		},
		{
			name: "トークンない",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, nil, http.MethodGet, "/notification/get", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}
