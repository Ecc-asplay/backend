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

func createTestUser(t *testing.T, userData CreateUserRequest) UserRsp {
	if userData.Username == "" && userData.Password == "" && userData.Email == "" {
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

	recorder := APITestBeforeLogin(t, userData, http.MethodPost, "/users")
	require.Equal(t, http.StatusOK, recorder.Code)

	var createdUser UserRsp
	err := json.Unmarshal(recorder.Body.Bytes(), &createdUser)
	require.NoError(t, err)

	fmt.Println(" ")
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
			name: "無効なデータ（メール）",
			body: CreateUserRequest{
				Username: gofakeit.Name(),
				Email:    "",
				Birth: pgtype.Date{
					Time:  util.RandomDate(),
					Valid: true,
				},
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
			recorder := APITestBeforeLogin(t, tc.body, http.MethodPost, "/users")
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteUserAPI(t *testing.T) {
	userDelData := createTestUser(t, CreateUserRequest{})
	token := "Bearer " + userDelData.Access_Token
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
			name: "トークンない",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, nil, http.MethodDelete, "/users/del", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestResetPasswordAPI(t *testing.T) {
	userRPData := createTestUser(t, CreateUserRequest{})
	token := "Bearer " + userRPData.Access_Token
	testCases := []struct {
		name          string
		token         string
		body          NewPasswordRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: NewPasswordRequest{
				NewPassword: util.RandomString(20),
			},
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: NewPasswordRequest{
				NewPassword: util.RandomString(20),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "新しいパスワードない",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/users/password", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestUpdateDiseaseConditionAPI(t *testing.T) {
	userData := createTestUser(t, CreateUserRequest{})
	token := "Bearer " + userData.Access_Token
	testCases := []struct {
		name          string
		body          NewDiseaseAndConditionRequest
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body: NewDiseaseAndConditionRequest{
				Disease:   util.RandomDisease(),
				Condition: util.RandomCondition(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: NewDiseaseAndConditionRequest{
				Disease:   util.RandomDisease(),
				Condition: util.RandomCondition(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "病歴と病状ない",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/users/disease-condition", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

// func TestUpdateEmailAPI(t *testing.T) {
// 	userData := createTestUser(t, CreateUserRequest{})
// 	token := "Bearer " + userData.Access_Token
// 	testCases := []struct {
// 		name          string
// 		token         string
// 		body          NewEmailRequest
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:  "OK",
// 			token: token,
// 			body: NewEmailRequest{
// 				NewEmail: gofakeit.Email(),
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "トークンない",
// 			body: NewEmailRequest{
// 				NewEmail: gofakeit.Email(),
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name:  "新しいメールない",
// 			token: token,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/users/email", tc.token)
// 			tc.checkResponse(recorder)
// 			fmt.Println(" ")
// 		})
// 	}
// }

// func TestUpdateIsPrivacyAPI(t *testing.T) {
// 	userData := createTestUser(t, CreateUserRequest{})
// 	token := "Bearer " + userData.Access_Token
// 	testCases := []struct {
// 		name          string
// 		token         string
// 		body          UpdatePrivacyRequest
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:  "OK",
// 			token: token,
// 			body: UpdatePrivacyRequest{
// 				IsPrivacy: !userData.User_Information.IsPrivacy,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "トークンない",
// 			body: UpdatePrivacyRequest{
// 				IsPrivacy: !userData.User_Information.IsPrivacy,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name:  "プライバシー更新ない",
// 			token: token,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/users/privacy", tc.token)
// 			tc.checkResponse(recorder)
// 			fmt.Println(" ")
// 		})
// 	}
// }

// func TestUpdateNameAPI(t *testing.T) {
// 	userData := createTestUser(t, CreateUserRequest{})
// 	token := "Bearer " + userData.Access_Token
// 	testCases := []struct {
// 		name          string
// 		token         string
// 		body          NewUsernameRequest
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:  "OK",
// 			token: token,
// 			body: NewUsernameRequest{
// 				NewUsername: gofakeit.Name(),
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "トークンない",
// 			body: NewUsernameRequest{
// 				NewUsername: gofakeit.Name(),
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name:  "新しい名前ない",
// 			token: token,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/users/name", tc.token)
// 			tc.checkResponse(recorder)
// 			fmt.Println(" ")
// 		})
// 	}
// }

func TestLoginUser(t *testing.T) {
	email := gofakeit.Email()
	userLoginData := createTestUser(t, CreateUserRequest{
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
				Email:    userLoginData.User_Information.Email,
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
				Email:    "asdadsadas",
				Password: "123qwecc",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "パスワードない",
			body: LoginRequest{
				Email: userLoginData.User_Information.Email,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestBeforeLogin(t, tc.body, http.MethodPost, "/login")
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
