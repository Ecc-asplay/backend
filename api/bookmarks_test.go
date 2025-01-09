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

func RandomCreateBookmarkAPI(t *testing.T, user1 UserRsp) db.Bookmark {
	token := "Bearer " + user1.Access_Token

	user2 := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

	post := RandomCreatePostAPI(t, user2)

	BMData := bookmarkRequest{
		PostID: post.PostID,
	}

	var createBM db.Bookmark

	t.Run("RandomBookmark", func(t *testing.T) {
		recorder := APITestAfterLogin(t, BMData, http.MethodPost, "/bookmark/add", token)
		bm, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		err = json.Unmarshal(bm, &createBM)
		require.NoError(t, err)

		fmt.Println(" ")
	})

	return createBM
}

func TestCreateBookmarkAPI(t *testing.T) {
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
	post := RandomCreatePostAPI(t, user)
	BMData := bookmarkRequest{
		PostID: post.PostID,
	}
	testCases := []struct {
		name          string
		token         string
		body          bookmarkRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  BMData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchBookmark(t, recorder.Body, BMData)
			},
		},
		{
			name: "トークンない",
			body: BMData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿ない",
			token: token,
			body: bookmarkRequest{
				PostID: util.CreateUUID(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPost, "/bookmark/add", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestGetBookmark(t *testing.T) {
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
		RandomCreateBookmarkAPI(t, user)
	}
	testCases := []struct {
		name          string
		token         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
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
			recorder := httptest.NewRecorder()
			server := newTestServer(t)
			require.NotEmpty(t, server)

			request, err := http.NewRequest(http.MethodGet, "/bookmark/get", nil)
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

func DeleteBookmark(t *testing.T) {
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
	bm := RandomCreateBookmarkAPI(t, user)

	testCases := []struct {
		name          string
		token         string
		body          bookmarkRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body: bookmarkRequest{
				PostID: bm.PostID,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: bookmarkRequest{
				PostID: bm.PostID,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "投稿IDない",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodDelete, "/bookmark/del", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func requireBodyMatchBookmark(t *testing.T, body *bytes.Buffer, bm bookmarkRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var getBookmark db.Bookmark
	err = json.Unmarshal(data, &getBookmark)
	require.NoError(t, err)

	require.NotEmpty(t, getBookmark.UserID)
	require.NotEmpty(t, getBookmark.CreatedAt)
	require.Equal(t, getBookmark.PostID, bm.PostID)
}
