package db

import (
	"context"
	"testing"

	"github.com/Ecc-asplay/backend/util"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func CreateRandomPost(t *testing.T, user User) Post {
	newPost := CreatePostParams{
		PostID:      util.CreateUUID(),
		UserID:      user.UserID,
		ShowID:      util.RandomString(10),
		Title:       gofakeit.BookTitle(),
		Feel:        util.RandomMood(),
		Content:     gofakeit.Sentence(30),
		Reaction:    rand.Int31(),
		Image:       gofakeit.ImagePng(200, 200),
		IsSensitive: util.RandomBool(),
		Status:      util.RandomStatus(),
	}

	post, err := testQueries.CreatePost(context.Background(), newPost)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, newPost.PostID, post.PostID)
	require.Equal(t, newPost.UserID, post.UserID)
	require.Equal(t, newPost.ShowID, post.ShowID)
	require.Equal(t, newPost.Title, post.Title)
	require.Equal(t, newPost.Feel, post.Feel)
	require.Equal(t, newPost.Content, post.Content)
	require.Equal(t, newPost.Reaction, post.Reaction)
	require.Equal(t, newPost.Image, post.Image)
	require.Equal(t, newPost.IsSensitive, post.IsSensitive)
	require.Equal(t, newPost.Status, post.Status)
	require.NotEmpty(t, post.CreatedAt.Time)

	return post
}

func TestCreatePost(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomPost(t, user)
}

func TestDeletePost(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	delete := DeletePostParams{
		UserID: user.UserID,
		PostID: post.PostID,
	}

	err := testQueries.DeletePost(context.Background(), delete)
	require.NoError(t, err)
}

func TestGetPostOfKeywords(t *testing.T) {
	user := CreateRandomUser(t)
	newPost := CreatePostParams{
		PostID:      util.CreateUUID(),
		UserID:      user.UserID,
		ShowID:      util.RandomString(10),
		Title:       "aaaaaaaaaaaaaaaaaaaa",
		Feel:        util.RandomMood(),
		Content:     "asdhavsjhsavdjah",
		Reaction:    rand.Int31(),
		Image:       gofakeit.ImagePng(200, 200),
		IsSensitive: util.RandomBool(),
		Status:      util.RandomStatus(),
	}

	post, err := testQueries.CreatePost(context.Background(), newPost)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	foundPost, err := testQueries.GetPostOfKeywords(context.Background(), "a")
	require.NoError(t, err)
	require.NotEmpty(t, foundPost)
}

func TestGetPostsList(t *testing.T) {

	for i := 0; i < 10; i++ {
		user := CreateRandomUser(t)
		for i := 0; i < 2; i++ {
			CreateRandomPost(t, user)
		}
	}

	postsList, err := testQueries.GetPostsList(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, postsList)
	require.GreaterOrEqual(t, len(postsList), 10)
}

func TestGetUserAllPosts(t *testing.T) {
	user := CreateRandomUser(t)
	for i := 0; i < 20; i++ {
		CreateRandomPost(t, user)
	}

	postsUserList, err := testQueries.GetUserAllPosts(context.Background(), user.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, postsUserList)
	require.GreaterOrEqual(t, len(postsUserList), 20)
}

func TestUpdatePosts(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	updateData := UpdatePostsParams{
		UserID:      user.UserID,
		PostID:      post.PostID,
		ShowID:      util.RandomString(10),
		Title:       gofakeit.BookTitle(),
		Feel:        util.RandomMood(),
		Content:     gofakeit.Sentence(30),
		Reaction:    rand.Int31(),
		Image:       gofakeit.ImagePng(300, 200),
		IsSensitive: util.RandomBool(),
	}

	newPost, err := testQueries.UpdatePosts(context.Background(), updateData)
	require.NoError(t, err)
	require.NotEmpty(t, newPost)
	require.Equal(t, updateData.PostID, newPost.PostID)
	require.Equal(t, updateData.UserID, newPost.UserID)
}
