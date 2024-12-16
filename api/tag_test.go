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

func RandomCreateTagAPI(t *testing.T) db.Tag {
	user := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + user.Access_Token
	post := RandomCreatePostAPI(t, user)

	tagData := CreateTagRequest{
		PostID:      post.PostID,
		TagComments: util.RandomString(5),
	}

	var createTag db.Tag

	t.Run("RandomTag", func(t *testing.T) {
		recoder := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(tagData)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/tag/add", bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", token)

		server.router.ServeHTTP(recoder, request)
		require.NotEmpty(t, recoder)

		tag, err := io.ReadAll(recoder.Body)
		require.NoError(t, err)

		err = json.Unmarshal(tag, &createTag)
		require.NoError(t, err)
		fmt.Println(" ")
	})

	return createTag
}

func TestCreateTagAPI(t *testing.T) {
	user := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + user.Access_Token
	post := RandomCreatePostAPI(t, user)

	tagData := CreateTagRequest{
		PostID:      post.PostID,
		TagComments: util.RandomString(5),
	}

	testCases := []struct {
		name          string
		token         string
		body          CreateTagRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  tagData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTag(t, recorder.Body, tagData)
			},
		},
		{
			name: "トークンない",
			body: tagData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			body: CreateTagRequest{
				TagComments: util.RandomString(5),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "タブない",
			token: token,
			body: CreateTagRequest{
				PostID: post.PostID,
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

			request, err := http.NewRequest(http.MethodPost, "/tag/add", bytes.NewReader(data))
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

func requireBodyMatchTag(t *testing.T, body *bytes.Buffer, tag CreateTagRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getTag db.Tag
	err = json.Unmarshal(data, &getTag)
	require.NoError(t, err)

	require.Equal(t, tag.PostID, getTag.PostID)
	require.Equal(t, tag.TagComments, getTag.TagComments)
}
