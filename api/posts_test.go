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

func RandomCreatePostAPI(t *testing.T, user UserRsp) db.Post {
	token := "Bearer " + user.Access_Token
	jsonData := ReturnContext()

	postData := CreatePostRequest{
		ShowID:  util.RandomString(10),
		Title:   "aaa",
		Feel:    util.RandomMood(),
		Content: jsonData,
		Status:  util.RandomStatus(),
	}

	var createPost db.Post

	t.Run("RandomPost", func(t *testing.T) {
		recorder := APITestAfterLogin(t, postData, http.MethodPost, "/post/add", token)
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
	jsonData := ReturnContext()

	postData := CreatePostRequest{
		ShowID:  util.RandomString(10),
		Title:   gofakeit.BookTitle(),
		Feel:    util.RandomMood(),
		Content: jsonData,
		Status:  util.RandomStatus(),
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
				ShowID:  util.RandomString(10),
				Title:   gofakeit.BookTitle(),
				Feel:    util.RandomMood(),
				Content: jsonData,
				Status:  util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "状態ない",
			token: token,
			body: CreatePostRequest{
				ShowID:  util.RandomString(10),
				Title:   gofakeit.BookTitle(),
				Feel:    util.RandomMood(),
				Content: jsonData,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPost, "/post/add", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestGetUserPostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

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
			recorder := APITestAfterLogin(t, nil, http.MethodGet, "/post/get", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestDeletPostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

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
			recorder := APITestAfterLogin(t, tc.body, http.MethodDelete, "/post/del", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestUpdatePostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

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
				PostID:  Post.PostID,
				ShowID:  util.RandomString(10),
				Title:   gofakeit.BookTitle(),
				Feel:    util.RandomMood(),
				Content: jsonData,
				Status:  util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: UpdatePostsRequest{
				PostID:  Post.PostID,
				ShowID:  util.RandomString(10),
				Title:   gofakeit.BookTitle(),
				Feel:    util.RandomMood(),
				Content: jsonData,
				Status:  util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			body: UpdatePostsRequest{
				ShowID:  util.RandomString(10),
				Title:   gofakeit.BookTitle(),
				Feel:    util.RandomMood(),
				Content: jsonData,
				Status:  util.RandomStatus(),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "状態ない",
			token: token,
			body: UpdatePostsRequest{
				PostID:  Post.PostID,
				ShowID:  util.RandomString(10),
				Title:   gofakeit.BookTitle(),
				Feel:    util.RandomMood(),
				Content: jsonData,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/post/update", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestSearchPostAPI(t *testing.T) {
	User := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})
	Post := RandomCreatePostAPI(t, User)

	testCases := []struct {
		name          string
		body          string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: Post.Title,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "見つけない",
			body: util.RandomString(10),
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestBeforeLogin(t, tc.body, http.MethodPost, "/post/search")
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
	require.False(t, getPost.IsSensitive)
	require.NotEmpty(t, getPost.CreatedAt)

	require.Equal(t, getPost.ShowID, post.ShowID)
	require.Equal(t, getPost.Title, post.Title)
	require.Equal(t, getPost.Feel, post.Feel)
	require.Equal(t, getPost.Status, post.Status)
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
