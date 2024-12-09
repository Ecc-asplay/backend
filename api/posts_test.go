package api

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Ecc-asplay/backend/util"
)

func RandomCreatePostData() CreatePostRequest {
	jsonData, err := gofakeit.JSON(nil)
	if err != nil {
		fmt.Println("Error generating JSON:", err)
	}

	return CreatePostRequest{
		ShowID:   util.RandomString(10),
		Title:    gofakeit.BookTitle(),
		Feel:     util.RandomMood(),
		Content:  jsonData,
		Reaction: rand.Int31(),
		Status:   util.RandomStatus(),
	}
}
