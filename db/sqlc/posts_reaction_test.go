package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomPostsReaction(t *testing.T, user User, post Post) PostsReaction {
	data := CreatePostsReactionParams{
		UserID:           user.UserID,
		PostID:           post.PostID,
		PReactionThanks:  util.RandomBool(),
		PReactionHelpful: util.RandomBool(),
		PReactionUseful:  util.RandomBool(),
		PReactionHeart:   util.RandomBool(),
	}

	reaction, err := testQueries.CreatePostsReaction(context.Background(), data)
	require.NoError(t, err)
	require.NotEmpty(t, reaction)
	require.Equal(t, reaction.UserID, user.UserID)
	require.Equal(t, reaction.PostID, post.PostID)
	require.Equal(t, reaction.PReactionThanks, data.PReactionThanks)
	require.Equal(t, reaction.PReactionHelpful, data.PReactionHelpful)
	require.Equal(t, reaction.PReactionUseful, data.PReactionUseful)
	require.Equal(t, reaction.PReactionHeart, data.PReactionHeart)
	require.NotEmpty(t, reaction.CreatedAt)
	return reaction
}

func TestCreateRandomPostsReaction(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	user2 := CreateRandomUser(t)
	CreateRandomPostsReaction(t, user2, post)
}

func TestGetPostsReaction(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomPostsReaction(t, user2, post)
	}

	allPostReaction, err := testQueries.GetPostsReaction(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotEmpty(t, allPostReaction)
	require.GreaterOrEqual(t, len(allPostReaction), 10)
}

func TestGetPostsHeartOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomPostsReaction(t, user2, post)
	}

	hearts, err := testQueries.GetPostsHeartOfTrue(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotZero(t, hearts)
}
func TestGetPostsHelpfulOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomPostsReaction(t, user2, post)
	}

	helpful, err := testQueries.GetPostsHelpfulOfTrue(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotZero(t, helpful)
}
func TestGetPostsThanksOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomPostsReaction(t, user2, post)
	}

	thanks, err := testQueries.GetPostsThanksOfTrue(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotZero(t, thanks)
}
func TestGetPostsUsefulOfTrue(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomPostsReaction(t, user2, post)
	}

	useful, err := testQueries.GetPostsUsefulOfTrue(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotZero(t, useful)
}

func TestUpdatePostsReactionHeart(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	user2 := CreateRandomUser(t)
	reaction := CreateRandomPostsReaction(t, user2, post)

	data := UpdatePostsReactionHeartParams{
		UserID: user2.UserID,
		PostID: reaction.PostID,
	}

	boo, err := testQueries.UpdatePostsReactionHeart(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, reaction.PReactionHeart)
}
func TestUpdatePostsReactionHelpful(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	user2 := CreateRandomUser(t)
	reaction := CreateRandomPostsReaction(t, user2, post)

	data := UpdatePostsReactionHelpfulParams{
		UserID: user2.UserID,
		PostID: reaction.PostID,
	}

	boo, err := testQueries.UpdatePostsReactionHelpful(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, reaction.PReactionHelpful)
}
func TestUpdatePostsReactionThanks(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	user2 := CreateRandomUser(t)
	reaction := CreateRandomPostsReaction(t, user2, post)

	data := UpdatePostsReactionThanksParams{
		UserID: user2.UserID,
		PostID: reaction.PostID,
	}

	boo, err := testQueries.UpdatePostsReactionThanks(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, reaction.PReactionHelpful)
}
func TestUpdatePostsReactionUseful(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	user2 := CreateRandomUser(t)
	reaction := CreateRandomPostsReaction(t, user2, post)

	data := UpdatePostsReactionUsefulParams{
		UserID: user2.UserID,
		PostID: reaction.PostID,
	}

	boo, err := testQueries.UpdatePostsReactionUseful(context.Background(), data)
	require.NoError(t, err)
	require.NotEqual(t, boo, reaction.PReactionHelpful)
}

func TestDeletePostsReaction(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	for i := 0; i < 10; i++ {
		user2 := CreateRandomUser(t)
		CreateRandomPostsReaction(t, user2, post)
	}

	allPostReaction, err := testQueries.GetPostsReaction(context.Background(), post.PostID)
	require.NoError(t, err)

	for _, reaction := range allPostReaction {
		if reaction.PReactionHeart == false && reaction.PReactionHelpful == false && reaction.PReactionThanks == false && reaction.PReactionUseful == false {
			data := DeletePostsReactionParams{
				UserID: reaction.UserID,
				PostID: reaction.PostID,
			}
			err = testQueries.DeletePostsReaction(context.Background(), data)
			require.NoError(t, err)
		}
	}
}
