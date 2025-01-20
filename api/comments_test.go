package api

import (
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

func RandomCreateCommentAPI(t *testing.T, postUser, user2 UserRsp) db.Comment {
	post := RandomCreatePostAPI(t, postUser)

	token := "Bearer " + user2.Access_Token
	CData := createCommentRequest{
		PostID:     post.PostID,
		Comments:   "aaaaaabbcdbdbhcjdsbhdjs",
		IsPublic:   false,
		IsCensored: false,
	}

	var createC db.Comment

	t.Run("RandomComment", func(t *testing.T) {
		recorder := APITestAfterLogin(t, CData, http.MethodPost, "/comment/create", token)
		c, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		err = json.Unmarshal(c, &createC)
		require.NoError(t, err)

		fmt.Println(" ")
	})

	return createC
}

func TestCreateCommentsAPI(t *testing.T) {
	postUser := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

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
	comment := RandomCreateCommentAPI(t, postUser, user2)

	token := "Bearer " + user2.Access_Token

	CData := createCommentRequest{
		PostID:     comment.PostID,
		Comments:   "aaaaaabbcdbdbhcjdsbhdjs",
		IsPublic:   false,
		IsCensored: false,
	}

	testCases := []struct {
		name          string
		token         string
		body          createCommentRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  CData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: CData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			body: createCommentRequest{
				Comments:   "aaaaaabbcdbdbhcjdsbhdjs",
				IsPublic:   false,
				IsCensored: false,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "コメントの内容ない",
			token: token,
			body: createCommentRequest{
				PostID:     comment.PostID,
				IsPublic:   false,
				IsCensored: false,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPost, "/comment/create", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestGetAllComments(t *testing.T) {
	postUser := RandomCreateUserAPI(t, CreateUserRequest{
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
		RandomCreateCommentAPI(t, postUser, user2)
	}

	token := "Bearer " + postUser.Access_Token

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
			recorder := APITestAfterLogin(t, nil, http.MethodGet, "/comment/all", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestGetPostCommentsListAPI(t *testing.T) {
	postUser := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

	var postid uuid.UUID
	for i := 0; i < 10; i++ {
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
		commentData := RandomCreateCommentAPI(t, postUser, user2)
		postid = commentData.PostID
	}

	token := "Bearer " + postUser.Access_Token

	testCases := []struct {
		name          string
		token         string
		body          uuid.UUID
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  postid,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: postid,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var url string
			if tc.body != uuid.Nil {
				url = `/comment/getlist/` + tc.body.String()
			} else {
				url = `/comment/getlist/`
			}
			recorder := APITestAfterLogin(t, nil, http.MethodGet, url, tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestUpdateCommentsAPI(t *testing.T) {
	postUser := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

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
	comment := RandomCreateCommentAPI(t, postUser, user2)
	token := "Bearer " + user2.Access_Token

	CData := UpdateCommentRequest{
		CommentID: comment.CommentID,
		Comments:  "aaaaaabbcdbdbhcjdsbhdjs",
		IsPublic:  !comment.IsPublic,
	}

	testCases := []struct {
		name          string
		token         string
		body          UpdateCommentRequest
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  CData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: CData,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			body: UpdateCommentRequest{
				Comments: "aaaaaabbcdbdbhcjdsbhdjs",
				IsPublic: !comment.IsPublic,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "コメントの内容ない",
			token: token,
			body: UpdateCommentRequest{
				CommentID: comment.CommentID,
				IsPublic:  !comment.IsPublic,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := APITestAfterLogin(t, tc.body, http.MethodPut, "/comment/update", tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}

func TestDeleteCommentsAPI(t *testing.T) {
	postUser := RandomCreateUserAPI(t, CreateUserRequest{
		Username: gofakeit.Name(),
		Email:    gofakeit.Email(),
		Birth: pgtype.Date{
			Time:  util.RandomDate(),
			Valid: true,
		},
		Gender:   util.RandomGender(),
		Password: util.RandomString(20),
	})

	var postid uuid.UUID
	for i := 0; i < 10; i++ {
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
		commentData := RandomCreateCommentAPI(t, postUser, user2)
		postid = commentData.PostID
	}

	token := "Bearer " + postUser.Access_Token

	testCases := []struct {
		name          string
		token         string
		body          uuid.UUID
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			token: token,
			body:  postid,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "トークンない",
			body: postid,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "投稿IDない",
			token: token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var url string
			if tc.body != uuid.Nil {
				url = `/comment/delete/` + tc.body.String()
			} else {
				url = `/comment/delete/`
			}
			recorder := APITestAfterLogin(t, nil, http.MethodDelete, url, tc.token)
			tc.checkResponse(recorder)
			fmt.Println(" ")
		})
	}
}
