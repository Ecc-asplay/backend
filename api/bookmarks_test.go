package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

func RandomCreateBookmarkAPI(t *testing.T, user1 UserRsp) db.Bookmark {
	token := "Bearer " + user1.Access_Token

	user2 := RandomCreateUserAPI(t, CreateUserRequest{})
	post := RandomCreatePostAPI(t, user2)

	BMData := bookmarkRequest{
		PostID: post.PostID,
	}

	var createBM db.Bookmark

	t.Run("RandomBookmark", func(t *testing.T) {
		recode := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(BMData)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/bookmark/add", bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", token)

		server.router.ServeHTTP(recode, request)
		require.NotEmpty(t, recode)

		bm, err := io.ReadAll(recode.Body)
		require.NoError(t, err)

		err = json.Unmarshal(bm, &createBM)
		require.NoError(t, err)

		fmt.Println(" ")
	})

	return createBM
}

func TestCreateBookmarkAPI(t *testing.T) {
	user := RandomCreateUserAPI(t, CreateUserRequest{})
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
			name:  "投稿IDない",
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
			recode := httptest.NewRecorder()
			server := newTestServer(t)
			require.NotEmpty(t, server)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/bookmark/add", bytes.NewReader(data))
			require.NoError(t, err)
			require.NotEmpty(t, request)

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", tc.token)

			server.router.ServeHTTP(recode, request)
			require.NotEmpty(t, recode)

			tc.checkResponse(recode)
			fmt.Println(" ")
		})
	}
}

// func DeleteBookmark(t *testing.T) {
// 	user := RandomCreateUserAPI(t, CreateUserRequest{})
// 	bm := wRandomCreateBookmarkAPI(t, user)

// }
