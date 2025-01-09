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

func RandomCreateBlockUser(t *testing.T, user1 UserRsp) db.Blockuser {
	token := "Bearer " + user1.Access_Token
	BlockUserAccount := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

	blockData := CreateBlockUserRequest{
		BlockUserID: BlockUserAccount.User_Information.UserID,
		Reason:      util.RandomString(10),
	}

	var createBlockUser db.Blockuser

	t.Run("RandomCreateBlockUser", func(t *testing.T) {
		recorder := APITestAfterLogin(t, blockData, http.MethodPost, "/block/create", token)

		buser, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		err = json.Unmarshal(buser, &createBlockUser)
		require.NoError(t, err)

		fmt.Println(" ")
	})
	return createBlockUser
}

func TestCreateBlockUser(t *testing.T) {
	user1 := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})
	token := "Bearer " + user1.Access_Token
	BlockUserAccount := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

	blockData := CreateBlockUserRequest{
		BlockUserID: BlockUserAccount.User_Information.UserID,
		Reason:      util.RandomString(10),
	}

	testCases := []struct {
		name          string
		token         string
		body          CreateBlockUserRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: token,
			body:  blockData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: CreateBlockUserRequest{
				BlockUserID: BlockUserAccount.User_Information.UserID,
				Reason:      util.RandomString(10),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ブロックユーザーIDない",
			body: CreateBlockUserRequest{
				Reason: util.RandomString(10),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "理由ない",
			body: CreateBlockUserRequest{
				BlockUserID: BlockUserAccount.User_Information.UserID,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPost, "/block/create", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestUnBlockUserAPI(t *testing.T) {
	user := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})
	token := "Bearer " + user.Access_Token
	blockuser := RandomCreateBlockUser(t, user)

	unblockData := UnblockUserRequest{
		BlockUserID: blockuser.BlockuserID,
	}

	testCases := []struct {
		name          string
		token         string
		body          UnblockUserRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: token,
			body:  unblockData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: unblockData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "ブロックユーザーIDない",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "ブロックユーザーID違った",
			token: token,
			body: UnblockUserRequest{
				BlockUserID: util.CreateUUID(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/block/update", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}

}

func TestGetBlockUserByUserAPI(t *testing.T) {
	user := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})
	token := "Bearer " + user.Access_Token
	for i := 0; i < 10; i++ {
		RandomCreateBlockUser(t, user)
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
			recorder := APITestAfterLogin(t, nil, http.MethodGet, "/block/get", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func requireBodyMatchBlockUser(t *testing.T, body *bytes.Buffer, buser CreateBlockUserRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var getBlockUser db.Blockuser
	err = json.Unmarshal(data, &getBlockUser)
	require.NoError(t, err)

	require.NotEmpty(t, getBlockUser.UserID)
	require.Equal(t, buser.BlockUserID, getBlockUser.BlockuserID)
	require.Equal(t, buser.Reason, getBlockUser.Reason)
	require.NotEmpty(t, getBlockUser.Status)
	require.NotEmpty(t, getBlockUser.BlockAt)
}
