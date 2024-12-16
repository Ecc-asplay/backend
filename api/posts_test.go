package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/util"
)

func RandomCreatePostAPI(t *testing.T, user UserRsp) db.Post {
	token := "Bearer " + user.Access_Token
	jsonData := ReturnContext()

	postData := CreatePostRequest{
		ShowID:   util.RandomString(10),
		Title:    gofakeit.BookTitle(),
		Feel:     util.RandomMood(),
		Content:  jsonData,
		Reaction: rand.Int31(),
		Status:   util.RandomStatus(),
	}

	var createPost db.Post

	t.Run("RandomPost", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(postData)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/post/add", bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", token)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)

		user, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		err = json.Unmarshal(user, &createPost)
		require.NoError(t, err)
		fmt.Println(" ")
	})

	return createPost
}

func TestCreatePostAPI(t *testing.T) {
	user := RandomCreateUserAPI(t, CreateUserRequest{})
	token := "Bearer " + user.Access_Token
	jsonData := ReturnContext()

	postData := CreatePostRequest{
		ShowID:   util.RandomString(10),
		Title:    gofakeit.BookTitle(),
		Feel:     util.RandomMood(),
		Content:  jsonData,
		Reaction: rand.Int31(),
		Status:   util.RandomStatus(),
	}

	testCases := []struct {
		name          string
		token         string
		body          CreatePostRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  postData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postData)
			},
		},
		{
			name: "トークンない",
			body: CreatePostRequest{
				ShowID:   util.RandomString(10),
				Title:    gofakeit.BookTitle(),
				Feel:     util.RandomMood(),
				Content:  jsonData,
				Reaction: rand.Int31(),
				Status:   util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "状態ない",
			token: token,
			body: CreatePostRequest{
				ShowID:   util.RandomString(10),
				Title:    gofakeit.BookTitle(),
				Feel:     util.RandomMood(),
				Content:  jsonData,
				Reaction: rand.Int31(),
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

			request, err := http.NewRequest(http.MethodPost, "/post/add", bytes.NewReader(data))
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

func TestGetUserPostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{})
	for i := 0; i < 10; i++ {
		RandomCreatePostAPI(t, User)
	}
	token := "Bearer " + User.Access_Token
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
			recorder := httptest.NewRecorder()
			server := newTestServer(t)
			require.NotEmpty(t, server)

			request, err := http.NewRequest(http.MethodGet, "/post/get", nil)
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

func TestDeletPostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{})
	Post := RandomCreatePostAPI(t, User)
	token := "Bearer " + User.Access_Token

	testCases := []struct {
		name          string
		token         string
		body          DeletePostRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body: DeletePostRequest{
				PostID: Post.PostID,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: DeletePostRequest{
				PostID: Post.PostID,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
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

			request, err := http.NewRequest(http.MethodDelete, "/post/del", bytes.NewReader(data))
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

func TestUpdatePostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{})
	Post := RandomCreatePostAPI(t, User)
	token := "Bearer " + User.Access_Token

	jsonData := ReturnContext()
	testCases := []struct {
		name          string
		token         string
		body          UpdatePostsRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body: UpdatePostsRequest{
				PostID:   Post.PostID,
				ShowID:   util.RandomString(10),
				Title:    gofakeit.BookTitle(),
				Feel:     util.RandomMood(),
				Content:  jsonData,
				Reaction: rand.Int31(),
				Status:   util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: UpdatePostsRequest{
				PostID:   Post.PostID,
				ShowID:   util.RandomString(10),
				Title:    gofakeit.BookTitle(),
				Feel:     util.RandomMood(),
				Content:  jsonData,
				Reaction: rand.Int31(),
				Status:   util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			body: UpdatePostsRequest{
				ShowID:   util.RandomString(10),
				Title:    gofakeit.BookTitle(),
				Feel:     util.RandomMood(),
				Content:  jsonData,
				Reaction: rand.Int31(),
				Status:   util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "状態ない",
			token: token,
			body: UpdatePostsRequest{
				PostID:   Post.PostID,
				ShowID:   util.RandomString(10),
				Title:    gofakeit.BookTitle(),
				Feel:     util.RandomMood(),
				Content:  jsonData,
				Reaction: rand.Int31(),
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

			request, err := http.NewRequest(http.MethodPut, "/post/update", bytes.NewReader(data))
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

func requireBodyMatchPost(t *testing.T, body *bytes.Buffer, post CreatePostRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getPost db.Post
	err = json.Unmarshal(data, &getPost)
	require.NoError(t, err)
	require.NotEmpty(t, getPost.UserID)
	require.NotEmpty(t, getPost.PostID)
	require.NotEmpty(t, getPost.UserID)
	require.NotEmpty(t, getPost.UserID)
	require.NotEmpty(t, getPost.UserID)
	require.NotEmpty(t, getPost.UserID)
	require.False(t, getPost.IsSensitive)
	require.NotEmpty(t, getPost.CreatedAt)
	require.Equal(t, post.ShowID, getPost.ShowID)
	require.Equal(t, post.Title, getPost.Title)
	require.Equal(t, post.Feel, getPost.Feel)
	require.NotEmpty(t, getPost.Content)
	require.Equal(t, post.Reaction, getPost.Reaction)
	require.Equal(t, post.Status, getPost.Status)
}

func ReturnContext() []byte {
	var Wsize []int
	for i := 0; i < 12; i++ {
		size := util.RandomInt(24) + 16
		Wsize = append(Wsize, size)
	}

	contentData := []map[string]interface{}{
		{
			"type": "paragraph",
			"children": []map[string]interface{}{
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[0],
					"text":     util.RandomHiragana(10),
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[1],
					"text":     util.RandomHiragana(10),
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[2],
					"text":     util.RandomHiragana(10),
				},
			},
		},
		{
			"type": "paragraph",
			"children": []map[string]interface{}{
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"italic":   util.RandomBool(),
					"fontsize": Wsize[3],
					"text":     util.RandomHiragana(10),
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"italic":   util.RandomBool(),
					"fontsize": Wsize[3],
					"text":     util.RandomHiragana(10),
				},
			},
		},
	}

	jsonData, err := json.MarshalIndent(contentData, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
	}

	return jsonData
}
