package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
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
		fmt.Println(" ")
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
			name: "ユーザーIDない",
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
			name: "無効なデータ（メール）",
			body: CreateUserRequset{
				Username: gofakeit.Name(),
				Email:    "",
				Birth:    util.RandomDate(),
				Gender:   util.RandomGender(),
				Password: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "無効なデータ（誕生日）",
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
			fmt.Println(" ")
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

func TestDeleteUserAPI(t *testing.T) {
	userData := RandomCreateUserAPI(t)
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			token:  token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "無効なデータ（ユーザーID）",
			userID: util.CreateUUID(),
			token:  token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "ユーザーIDない",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			token:  "",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			server := newTestServer(t)
			require.NotEmpty(t, server)

			data, err := json.Marshal(tc.userID)
			require.NoError(t, err)

			url := fmt.Sprintf("/users/%s", tc.userID)
			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", tc.token)

			server.router.ServeHTTP(recorder, request)
			require.NotEmpty(t, recorder)

			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestResetPasswordAPI(t *testing.T) {
	type newPw struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	userData := RandomCreateUserAPI(t)
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		newPassword   newPw
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			newPassword: newPw{
				NewPassword: util.RandomString(20),
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				// requireBodyMatchUser(t, recorder.Body,)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			newPassword: newPw{
				NewPassword: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "新しいパスワードない",
			userID: userData.User_Information.UserID,
			token:  token,
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

			data, err := json.Marshal(tc.newPassword)
			require.NoError(t, err)

			url := fmt.Sprintf("/users/%s/password", tc.userID)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", tc.token)

			server.router.ServeHTTP(recorder, request)
			require.NotEmpty(t, recorder)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestUpdatediseaseConditionAPI(t *testing.T) {
	type newDiseaseAndCondition struct {
		Disease   string `json:"disease" binding:"required"`
		Condition string `json:"condition" binding:"required"`
	}

	userData := RandomCreateUserAPI(t)
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name                   string
		userID                 uuid.UUID
		newDiseaseAndCondition newDiseaseAndCondition
		token                  string
		checkResponse          func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			newDiseaseAndCondition: newDiseaseAndCondition{
				Disease:   util.RandomDisease(),
				Condition: util.RandomCondition(),
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			newDiseaseAndCondition: newDiseaseAndCondition{
				Disease:   util.RandomDisease(),
				Condition: util.RandomCondition(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "病歴と病状ない",
			userID: userData.User_Information.UserID,
			token:  token,
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

			data, err := json.Marshal(tc.newDiseaseAndCondition)
			require.NoError(t, err)

			url := fmt.Sprintf("/users/%s/disease-condition", tc.userID)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", tc.token)

			server.router.ServeHTTP(recorder, request)
			require.NotEmpty(t, recorder)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestUpdateEmailAPI(t *testing.T) {
	type newEmail struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	userData := RandomCreateUserAPI(t)
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		newEmail      newEmail
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			newEmail: newEmail{
				NewEmail: gofakeit.Email(),
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			newEmail: newEmail{
				NewEmail: gofakeit.Email(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "メールない",
			userID: userData.User_Information.UserID,
			token:  token,
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

			data, err := json.Marshal(tc.newEmail)
			require.NoError(t, err)

			url := fmt.Sprintf("/users/%s/email", tc.userID)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", tc.token)

			server.router.ServeHTTP(recorder, request)
			require.NotEmpty(t, recorder)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}
