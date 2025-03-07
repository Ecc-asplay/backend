package db

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomPost(t *testing.T, user User) Post {
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
					"text":     "頭",
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[1],
					"text":     "が痛",
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[2],
					"text":     "い",
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
					"text":     "頭痛が",
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"italic":   util.RandomBool(),
					"fontsize": Wsize[3],
					"text":     "痛い",
				},
			},
		},
	}
	contentJson, err := json.MarshalIndent(contentData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	newPost := CreatePostParams{
		PostID:      util.CreateUUID(),
		UserID:      user.UserID,
		ShowID:      util.RandomString(10),
		Title:       gofakeit.BookTitle(),
		Feel:        util.RandomMood(),
		Content:     contentJson,
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
					"strike":   util.RandomBool(),
					"fontsize": Wsize[5],
					"text":     "四十",
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"strike":   util.RandomBool(),
					"fontsize": Wsize[6],
					"text":     "肩",
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"strike":   util.RandomBool(),
					"fontsize": Wsize[7],
					"text":     "がつらい",
				},
			},
		},
		{
			"type": "paragraph",
			"children": []map[string]interface{}{
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[8],
					"text":     "深爪",
				},
				{
					"bold":     util.RandomBool(),
					"color":    gofakeit.Color(),
					"fontsize": Wsize[9],
					"text":     "が深い",
				},
			},
		},
	}
	contentJson, err := json.MarshalIndent(contentData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	updateData := UpdatePostsParams{
		UserID:      user.UserID,
		PostID:      post.PostID,
		ShowID:      util.RandomString(10),
		Title:       gofakeit.BookTitle(),
		Feel:        util.RandomMood(),
		Content:     contentJson,
		IsSensitive: util.RandomBool(),
		Status:      util.RandomStatus(),
	}

	newPost, err := testQueries.UpdatePosts(context.Background(), updateData)
	require.NoError(t, err)
	require.NotEmpty(t, newPost)
	require.Equal(t, updateData.PostID, newPost.PostID)
	require.Equal(t, updateData.UserID, newPost.UserID)
}

func TestGetOnePost(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	p, err := testQueries.GetOnePost(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotEmpty(t, p)
}
