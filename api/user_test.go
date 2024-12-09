package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

type UserRsp struct {
	Access_Token     string
	TakeAt           pgtype.Timestamp
	User_Information db.User
}

func RandomCreateUserAPI(t *testing.T) UserRsp {
	userData := CreateUserRequset{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth:    util.RandomDate(),
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	}
	var createdUser UserRsp

	t.Run("RandomCreate", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(userData)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)

		user, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)
		err = json.Unmarshal(user, &createdUser)
		require.NoError(t, err)
	})
	return createdUser
}

func TestCreateUserAPI(t *testing.T) {
	userData := CreateUserRequset{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth:    util.RandomDate(),
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	}

	testCases := []struct {
		name          string
		body          CreateUserRequset
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: userData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, userData)
			},
		},
		{
			name: "Missing Username",
			body: CreateUserRequset{
				Email:    gofakeit.Email(),
				Birth:    util.RandomDate(),
				Gender:   util.RandomGender(),
				Password: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: CreateUserRequset{
				Username: gofakeit.Name(),
				Email:    "invalid-email",
				Birth:    util.RandomDate(),
				Gender:   util.RandomGender(),
				Password: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid Birth Date",
			body: CreateUserRequset{
				Username: gofakeit.Name(),
				Email:    gofakeit.Email(),
				Birth:    time.Time{},
				Gender:   util.RandomGender(),
				Password: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			server := newTestServer(t)
			require.NotEmpty(t, server)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			server.router.ServeHTTP(recorder, request)
			require.NotEmpty(t, recorder)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user CreateUserRequset) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotUser UserRsp
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.User_Information.Username)
	require.Equal(t, user.Email, gotUser.User_Information.Email)
	require.Equal(t, user.Gender, gotUser.User_Information.Gender)
	require.NotEmpty(t, gotUser.User_Information.Hashpassword)
}

// func TestDeleteUserAPI(t *testing.T) {
// 	userData := RandomCreateUserAPI(t)
// 	testCases := []struct {
// 		name          string
// 		userID        uuid.UUID
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:   "OK",
// 			userID: userData.User_Information.UserID,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name:   "Invalid UserID",
// 			userID: util.CreateUUID(),
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "missing UserID",
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			recorder := httptest.NewRecorder()
// 			server := newTestServer(t)
// 			require.NotEmpty(t, server)

// 			data, err := json.Marshal(tc.userID)
// 			require.NoError(t, err)

// 			url := fmt.Sprintf("/users/%s", tc.userID)
// 			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
// 			require.NoError(t, err)
// 			require.NotEmpty(t, request)

// 			token := "Bearer " + userData.Access_Token
// 			// fmt.Println("------token--------", token)
// 			request.Header.Set("Content-Type", "application/json")
// 			request.Header.Set("Authorization", token)

// 			server.router.ServeHTTP(recorder, request)
// 			require.NotEmpty(t, recorder)

// 			tc.checkResponse(recorder)
// 			fmt.Println(" ")
// 		})
// 	}
// }
