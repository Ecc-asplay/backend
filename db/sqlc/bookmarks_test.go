package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomBookmark(t *testing.T, user User, post Post) Bookmark {

	bookmarkData := CreateBookmarksParams{
		UserID: user.UserID,
		PostID: post.PostID,
	}
	bookmark, err := testQueries.CreateBookmarks(context.Background(), bookmarkData)
	require.NoError(t, err)
	require.NotEmpty(t, bookmark)
	require.Equal(t, bookmark.UserID, user.UserID)
	require.Equal(t, bookmark.PostID, post.PostID)
	require.NotEmpty(t, bookmark.CreatedAt)

	return bookmark
}

func TestCreateBookmark(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)
	post := CreateRandomPost(t, user2)
	CreateRandomBookmark(t, user1, post)
}

func TestDeleteBookmarks(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)
	post := CreateRandomPost(t, user2)
	bm := CreateRandomBookmark(t, user1, post)

	delete := DeleteBookmarksParams{
		UserID: user1.UserID,
		PostID: bm.PostID,
	}

	err := testQueries.DeleteBookmarks(context.Background(), delete)
	require.NoError(t, err)
}

func TestGetAllBookmarks(t *testing.T) {
	user1 := CreateRandomUser(t)
	for i := 0; i < 20; i++ {
		user2 := CreateRandomUser(t)
		post := CreateRandomPost(t, user2)
		CreateRandomBookmark(t, user1, post)
	}

	allBm, err := testQueries.GetAllBookmarks(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, allBm)
	require.GreaterOrEqual(t, len(allBm), 1)
}
