package db

import (
	"testing"

	"github.com/Ecc-asplay/backend/util"
	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/exp/rand"
)

func CreateRandomPost(t *testing.T) {
	user := CreateRandomUser(t)

	newPost := CreatePostParams{
		PostID:      util.CreateUUID(),
		UserID:      user.UserID,
		ShowID:      util.RandomString(10),
		Title:       gofakeit.BookTitle(),
		Feel:        util.RandomMood(),
		Content:     gofakeit.Sentence(30),
		Reaction:    rand.Int31(),
		Image:       []byte{gofakeit.ImagePng(200, 200)},
		IsSensitive: gofakeit.Bool(),
	}
}
