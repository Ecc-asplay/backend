package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomTag(t *testing.T, user User, post Post) Tag {

	tagData := CreateTagParams{
		PostID:      post.PostID,
		TagComments: util.RandomString(10),
	}
	tag, err := testQueries.CreateTag(context.Background(), tagData)
	require.NoError(t, err)
	require.NotEmpty(t, tag)
	require.Equal(t, tag.PostID, post.PostID)
	require.Equal(t, tag.TagComments, tagData.TagComments)

	return tag
}

func TestCreateRandomTag(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	CreateRandomTag(t, user, post)
}

func TestGetTag(t *testing.T) {
	/// ----------------------------------------------------------------
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)

	tagData := CreateTagParams{
		PostID:      post.PostID,
		TagComments: "yyyy",
	}

	newtag, err := testQueries.CreateTag(context.Background(), tagData)
	require.NoError(t, err)
	require.NotEmpty(t, newtag)
	require.Equal(t, newtag.PostID, post.PostID)
	require.Equal(t, newtag.TagComments, tagData.TagComments)

	foundTag, err := testQueries.GetTag(context.Background(), "y")
	require.NoError(t, err)
	require.NotEmpty(t, foundTag)
	require.Greater(t, len(foundTag), 0)
}
