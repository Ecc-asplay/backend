package db

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

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
		{
			"type": "paragraph",
			"children": []map[string]interface{}{
				{
					"color":     gofakeit.Color(),
					"underline": util.RandomBool(),
					"fontsize":  Wsize[4],
					"text":      "偏頭痛が痛い",
				},
			},
		},
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

	img := [][]byte{
		gofakeit.ImageJpeg(100, 100),
		gofakeit.ImageJpeg(100, 100),
		gofakeit.ImageJpeg(100, 100),
		gofakeit.ImageJpeg(100, 100),
		gofakeit.ImageJpeg(100, 100),
	}

	// imageData 數據結構
	imageData := []map[string]interface{}{
		{
			"page":   1,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
		{
			"page":   2,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
		{
			"page":   3,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
	}

	// 將 imageData 序列化為 JSON 格式
	imageJson, err := json.MarshalIndent(imageData, "", "  ")
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
		Images:      imageJson,
		Reaction:    rand.Int31(),
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
	require.Equal(t, newPost.Reaction, post.Reaction)
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
		{
			"type": "paragraph",
			"children": []map[string]interface{}{
				{
					"color":     gofakeit.Color(),
					"underline": util.RandomBool(),
					"fontsize":  Wsize[4],
					"text":      "偏頭痛が痛い",
				},
			},
		},
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

	img := [][]byte{gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200)}
	imageData := []map[string]interface{}{
		{
			"page":   1,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
		{
			"page":   2,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
		{
			"page":   3,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
	}
	imageJson, err := json.MarshalIndent(imageData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	newPost := CreatePostParams{
		PostID:      util.CreateUUID(),
		UserID:      user.UserID,
		ShowID:      util.RandomString(10),
		Title:       "aaaaaaaaaaaaaaaaaaaa",
		Feel:        util.RandomMood(),
		Content:     contentJson,
		Images:      imageJson,
		Reaction:    rand.Int31(),
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
		{
			"type": "paragraph",
			"children": []map[string]interface{}{
				{
					"color":     gofakeit.Color(),
					"underline": util.RandomBool(),
					"fontsize":  Wsize[4],
					"text":      "偏頭痛が痛い",
				},
			},
		},
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
	img := [][]byte{gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200), gofakeit.ImageJpeg(200, 200)}
	imageData := []map[string]interface{}{
		{
			"page":   1,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
		{
			"page":   2,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
		{
			"page":   3,
			"image1": img[0],
			"image2": img[1],
			"image3": img[2],
			"image4": img[3],
			"image5": img[4],
		},
	}
	imageJson, err := json.MarshalIndent(imageData, "", "  ")
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
		Images:      imageJson,
		Reaction:    rand.Int31(),
		IsSensitive: util.RandomBool(),
	}

	newPost, err := testQueries.UpdatePosts(context.Background(), updateData)
	require.NoError(t, err)
	require.NotEmpty(t, newPost)
	require.Equal(t, updateData.PostID, newPost.PostID)
	require.Equal(t, updateData.UserID, newPost.UserID)
}
