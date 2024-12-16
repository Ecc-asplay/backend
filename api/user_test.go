package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func RandomCreateUserAPI(t *testing.T, CusData CreateUserRequest) UserRsp {
	var userData CreateUserRequest
	if CusData.Email != "" && CusData.Username != "" && CusData.Password != "" {
		userData = CusData
	} else {
		userData = CreateUserRequest{
			Username: gofakeit.Name(),
			Email:    gofakeit.Email(),
			Birth: pgtype.Date{
				Time:  util.RandomDate(),
				Valid: true,
			},
			Gender:   util.RandomGender(),
			Password: util.RandomString(20),
		}
	}

	var createdUser UserRsp

	t.Run("RandomCreate", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(userData)
		require.NoError(t, err)
		require.NotEmpty(t, data)

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
	userData := CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	}

	testCases := []struct {
		name          string
		body          CreateUserRequest
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
			body: CreateUserRequest{
				Email: gofakeit.Email(),
				Birth: pgtype.Date{
					Time:  util.RandomDate(),
					Valid: true,
				}, Gender: util.RandomGender(),
				Password: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "無効なデータ（メール）",
			body: CreateUserRequest{
				Username: gofakeit.Name(),
				Email:    "",
				Birth: pgtype.Date{
					Time:  util.RandomDate(),
					Valid: true,
				}, Gender: util.RandomGender(),
				Password: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "無効なデータ（誕生日）",
			body: CreateUserRequest{
				Username: gofakeit.Name(),
				Email:    gofakeit.Email(),
				Birth:    pgtype.Date{},
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
			require.NotEmpty(t, data)

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

func TestDeleteUserAPI(t *testing.T) {
	userData := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		body          uuid.UUID
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			body:  userData.User_Information.UserID,
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "無効なデータ（ユーザーID）",
			body:  util.CreateUUID(),
			token: token,
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
			name:  "トークンない",
			body:  userData.User_Information.UserID,
			token: "",
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodDelete, "/users/del", bytes.NewReader(data))
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
	userData := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		body          NewPasswordRequest
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			body: NewPasswordRequest{
				NewPassword: util.RandomString(20),
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			body: NewPasswordRequest{
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/users/password", bytes.NewReader(data))
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
	userData := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		body          NewDiseaseAndConditionRequest
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			body: NewDiseaseAndConditionRequest{
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
			body: NewDiseaseAndConditionRequest{
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/users/disease-condition", bytes.NewReader(data))
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
	userData := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		body          NewEmailRequest
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			body: NewEmailRequest{
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
			body: NewEmailRequest{
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/users/email", bytes.NewReader(data))
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

func TestUpdateIsPrivacyAPI(t *testing.T) {
	type Privacy struct {
		IsPrivacy bool `json:"is_privacy" binding:"required"`
	}

	userData := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		body          Privacy
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			body: Privacy{
				IsPrivacy: !userData.User_Information.IsPrivacy,
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			body: Privacy{
				IsPrivacy: !userData.User_Information.IsPrivacy,
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/users/privacy", bytes.NewReader(data))
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

func TestUpdateNameAPI(t *testing.T) {
	userData := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		userID        uuid.UUID
		body          NewUsernameRequest
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: userData.User_Information.UserID,
			body: NewUsernameRequest{
				NewUsername: gofakeit.Name(),
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "トークンない",
			userID: userData.User_Information.UserID,
			body: NewUsernameRequest{
				NewUsername: gofakeit.Name(),
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, "/users/name", bytes.NewReader(data))
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

func TestLoginUser(t *testing.T) {
	email := gofakeit.Email()
	userData := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    email,
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   gofakeit.Gender(),
		Password: "123qwecc",
	})

	testCases := []struct {
		name          string
		body          LoginRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: LoginRequest{
				Email:    userData.User_Information.Email,
				Password: "123qwecc",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "メールない",
			body: LoginRequest{
				Password: "123qwecc",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "無効なデータ（メール&パスワード）",
			body: LoginRequest{
				Email:    gofakeit.Email(),
				Password: "123qwecc",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "パスワードない",
			body: LoginRequest{
				Email: userData.User_Information.Email,
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

			request, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			server.router.ServeHTTP(recorder, request)
			require.NotEmpty(t, recorder)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}

}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user CreateUserRequest) {
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
