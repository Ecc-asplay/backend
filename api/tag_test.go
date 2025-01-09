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

func RandomCreateTagAPI(t *testing.T) db.Tag {
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

	tagData := CreateTagRequest{
		PostID:      post.PostID,
		TagComments: util.RandomString(5),
	}

	var createTag db.Tag

	t.Run("RandomTag", func(t *testing.T) {
		recorder := APITestAfterLogin(t, tagData, http.MethodPost, "/tag/add", token)
		require.NotEmpty(t, recorder)

		tag, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		err = json.Unmarshal(tag, &createTag)
		require.NoError(t, err)
		fmt.Println(" ")
	})

	return createTag
}

func TestCreateTagAPI(t *testing.T) {
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
			recorder := APITestAfterLogin(t, tc.body, http.MethodPost, "/tag/add", tc.token)
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
