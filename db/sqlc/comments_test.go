package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomComment(t *testing.T, user User, post Post) Comment {
	commentData := CreateCommentsParams{
		CommentID:  util.CreateUUID(),
		UserID:     user.UserID,
		PostUser:   post.UserID,
		PostID:     post.PostID,
		Status:     util.RandomStatus(),
		IsPublic:   util.RandomBool(),
		Comments:   gofakeit.Sentence(20),
		IsCensored: util.RandomBool(),
	}

	comment, err := testQueries.CreateComments(context.Background(), commentData)
	require.NoError(t, err)
	require.NotEmpty(t, comment)
	require.Equal(t, comment.CommentID, commentData.CommentID)
	require.Equal(t, comment.UserID, commentData.UserID)
	require.Equal(t, comment.PostUser, commentData.PostUser)
	require.Equal(t, comment.PostID, commentData.PostID)
	require.Equal(t, comment.Status, commentData.Status)
	require.Equal(t, comment.IsPublic, commentData.IsPublic)
	require.Equal(t, comment.Comments, commentData.Comments)
	require.Equal(t, comment.IsCensored, commentData.IsCensored)

	return comment
}

func TestCreateComments(t *testing.T) {
	user1 := CreateRandomUser(t)
	post := CreateRandomPost(t, user1)
	user2 := CreateRandomUser(t)
	CreateRandomComment(t, user2, post)
}

func TestDeleteComments(t *testing.T) {
	user1 := CreateRandomUser(t)
	post := CreateRandomPost(t, user1)
	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	err := testQueries.DeleteComments(context.Background(), comment.CommentID)
	require.NoError(t, err)
}

func TestGetCommentsList(t *testing.T) {
	user1 := CreateRandomUser(t)
	post := CreateRandomPost(t, user1)
	for i := 0; i < 20; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomComment(t, user2, post)
	}

	allCommentData, err := testQueries.GetCommentsList(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotEmpty(t, allCommentData)
	require.GreaterOrEqual(t, len(allCommentData), 1)
}

func TestGetMyComment(t *testing.T) {
	user1 := CreateRandomUser(t)
	for i := 0; i < 20; i++ {
		user2 := CreateRandomUser(t)
		post := CreateRandomPost(t, user2)
		CreateRandomComment(t, user1, post)
	}

	allMyComments, err := testQueries.GetMyComments(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, allMyComments)
	require.GreaterOrEqual(t, len(allMyComments), 1)
}

func TestUpdateComments(t *testing.T) {
	user1 := CreateRandomUser(t)
	post := CreateRandomPost(t, user1)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	updateCommentData := UpdateCommentsParams{
		CommentID:  comment.CommentID,
		Status:     util.RandomStatus(),
		IsPublic:   util.RandomBool(),
		IsCensored: util.RandomBool(),
	}

	updatedComment, err := testQueries.UpdateComments(context.Background(), updateCommentData)
	require.NoError(t, err)
	require.NotEmpty(t, updatedComment)
	require.Equal(t, updatedComment.CommentID, comment.CommentID)
	require.Equal(t, updatedComment.UserID, comment.UserID)
	require.Equal(t, updatedComment.PostID, comment.PostID)
}

func TestGetAllComment(t *testing.T) {
	user1 := CreateRandomUser(t)
	post := CreateRandomPost(t, user1)
	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomComment(t, user2, post)
	}

	allCommentData, err := testQueries.GetAllComments(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, allCommentData)
	require.GreaterOrEqual(t, len(allCommentData), 1)
}

func TestGetPublicComments(t *testing.T) {
	user1 := CreateRandomUser(t)
	post := CreateRandomPost(t, user1)
	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomComment(t, user2, post)
	}

	publicComment, err := testQueries.GetAllPublicComments(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, publicComment)
	require.GreaterOrEqual(t, len(publicComment), 1)
}
