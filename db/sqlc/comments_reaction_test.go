package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomCommentsReaction(t *testing.T, user User, comment Comment) CommentsReaction {
	data := CreateCommentsReactionParams{
		UserID:           user.UserID,
		CommentID:        comment.CommentID,
		CReactionThanks:  util.RandomBool(),
		CReactionHelpful: util.RandomBool(),
		CReactionUseful:  util.RandomBool(),
		CReactionHeart:   util.RandomBool(),
	}

	reaction, err := testQueries.CreateCommentsReaction(context.Background(), data)
	require.NoError(t, err)
	require.NotEmpty(t, reaction)
	require.Equal(t, reaction.UserID, user.UserID)
	require.Equal(t, reaction.CommentID, comment.CommentID)
	require.Equal(t, reaction.CReactionThanks, data.CReactionThanks)
	require.Equal(t, reaction.CReactionHelpful, data.CReactionHelpful)
	require.Equal(t, reaction.CReactionUseful, data.CReactionUseful)
	require.Equal(t, reaction.CReactionHeart, data.CReactionHeart)
	require.NotEmpty(t, reaction.CreatedAt)
	return reaction
}

func TestCreateRandomCommentsReaction(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	user3 := CreateRandomUser(t)
	CreateRandomCommentsReaction(t, user3, comment)
}

func TestGetCommentsReaction(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)
	for i := 0; i < 10; i++ {
		user3 := CreateRandomUser(t)
		CreateRandomCommentsReaction(t, user3, comment)
	}

	allPostReaction, err := testQueries.GetCommentsReaction(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, allPostReaction)
	require.GreaterOrEqual(t, len(allPostReaction), 10)
}

func TestGetCommentsHeartOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)
	for i := 0; i < 10; i++ {
		user3 := CreateRandomUser(t)
		CreateRandomCommentsReaction(t, user3, comment)
	}

	c_hearts, err := testQueries.GetCommentsHeartOfTrue(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotZero(t, c_hearts)
}

func TestGetCommentsHelpfulOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)
	for i := 0; i < 10; i++ {
		user3 := CreateRandomUser(t)
		CreateRandomCommentsReaction(t, user3, comment)
	}

	c_helpful, err := testQueries.GetCommentsHelpfulOfTrue(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotZero(t, c_helpful)
}

func TestGetCommentsThanksOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)
	for i := 0; i < 10; i++ {
		user3 := CreateRandomUser(t)
		CreateRandomCommentsReaction(t, user3, comment)
	}

	c_thanks, err := testQueries.GetCommentsThanksOfTrue(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotZero(t, c_thanks)
}

func TestGetCommentsUsefulOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)
	for i := 0; i < 10; i++ {
		user3 := CreateRandomUser(t)
		CreateRandomCommentsReaction(t, user3, comment)
	}

	c_useful, err := testQueries.GetCommentsUsefulOfTrue(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotZero(t, c_useful)
}

func TestUpdateCommentsReactionThanks(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	user3 := CreateRandomUser(t)
	commentReaction := CreateRandomCommentsReaction(t, user3, comment)

	data := UpdateCommentsReactionThanksParams{
		UserID:    user3.UserID,
		CommentID: comment.CommentID,
	}

	boo, err := testQueries.UpdateCommentsReactionThanks(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, commentReaction.CReactionThanks)
}
func TestUpdateCommentsReactionHeart(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	user3 := CreateRandomUser(t)
	commentReaction := CreateRandomCommentsReaction(t, user3, comment)

	data := UpdateCommentsReactionHeartParams{
		UserID:    user3.UserID,
		CommentID: comment.CommentID,
	}

	boo, err := testQueries.UpdateCommentsReactionHeart(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, commentReaction.CReactionHeart)
}
func TestUpdateCommentsReactionHelpful(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	user3 := CreateRandomUser(t)
	commentReaction := CreateRandomCommentsReaction(t, user3, comment)

	data := UpdateCommentsReactionHelpfulParams{
		UserID:    user3.UserID,
		CommentID: comment.CommentID,
	}

	boo, err := testQueries.UpdateCommentsReactionHelpful(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, commentReaction.CReactionHelpful)
}
func TestUpdateCommentsReactionUseful(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)

	user3 := CreateRandomUser(t)
	commentReaction := CreateRandomCommentsReaction(t, user3, comment)

	data := UpdateCommentsReactionUsefulParams{
		UserID:    user3.UserID,
		CommentID: comment.CommentID,
	}

	boo, err := testQueries.UpdateCommentsReactionUseful(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, commentReaction.CReactionUseful)
}

func TestDeleteCommentsReaction(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	user2 := CreateRandomUser(t)
	comment := CreateRandomComment(t, user2, post)
	for i := 0; i < 10; i++ {
		user3 := CreateRandomUser(t)
		CreateRandomCommentsReaction(t, user3, comment)
	}

	allCommentsReaction, err := testQueries.GetCommentsReaction(context.Background(), comment.CommentID)
	require.NoError(t, err)
	require.NotEmpty(t, allCommentsReaction)

	for _, reaction := range allCommentsReaction {
		if reaction.CReactionHeart == false && reaction.CReactionHelpful == false && reaction.CReactionThanks == false && reaction.CReactionUseful == false {
			data := DeleteCommentsReactionParams{
				UserID:    reaction.UserID,
				CommentID: reaction.CommentID,
			}
			err = testQueries.DeleteCommentsReaction(context.Background(), data)
			require.NoError(t, err)
		}
	}
}
