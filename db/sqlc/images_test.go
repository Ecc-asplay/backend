package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func CreateRandomImage(t *testing.T, post Post) Image {
	creaImg := CreateImageParams{
		PostID: post.PostID,
		Page:   gofakeit.Int32(),
		Img1:   gofakeit.ImageJpeg(50, 50),
		Img2:   gofakeit.ImageJpeg(50, 50),
		Img3:   gofakeit.ImageJpeg(50, 50),
		Img4:   gofakeit.ImageJpeg(50, 50),
		Img5:   gofakeit.ImageJpeg(50, 50),
	}

	img, err := testQueries.CreateImage(context.Background(), creaImg)
	require.NoError(t, err)
	require.NotEmpty(t, img)
	require.Equal(t, img.PostID, creaImg.PostID)
	require.Equal(t, img.Page, creaImg.Page)
	require.Equal(t, img.Img1, creaImg.Img1)
	require.Equal(t, img.Img2, creaImg.Img2)
	require.Equal(t, img.Img3, creaImg.Img3)
	require.Equal(t, img.Img4, creaImg.Img4)
	require.Equal(t, img.Img5, creaImg.Img5)
	require.NotEmpty(t, img.CreatedAt)

	return img
}

func TestCreateImage(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	CreateRandomImage(t, post)
}

func TestGetImage(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	imgData := CreateRandomImage(t, post)

	img, err := testQueries.GetImage(context.Background(), post.PostID)
	require.NoError(t, err)
	require.NotEmpty(t, img)
	require.Equal(t, img[0].PostID, imgData.PostID)
	require.Equal(t, img[0].Page, imgData.Page)
	require.Equal(t, img[0].Img1, imgData.Img1)
	require.Equal(t, img[0].Img2, imgData.Img2)
	require.Equal(t, img[0].Img3, imgData.Img3)
	require.Equal(t, img[0].Img4, imgData.Img4)
	require.Equal(t, img[0].Img5, imgData.Img5)
	require.Equal(t, img[0].CreatedAt, imgData.CreatedAt)
}

func TestUpdateImages(t *testing.T) {
	user := CreateRandomUser(t)
	post := CreateRandomPost(t, user)
	imgData := CreateRandomImage(t, post)

	newImg := UpdateImageParams{
		PostID: post.PostID,
		Page:   gofakeit.Int32(),
		Img1:   gofakeit.ImageJpeg(51, 50),
		Img2:   gofakeit.ImageJpeg(51, 50),
		Img3:   gofakeit.ImageJpeg(51, 50),
		Img4:   gofakeit.ImageJpeg(51, 50),
		Img5:   gofakeit.ImageJpeg(51, 50),
	}

	newImgData, err := testQueries.UpdateImage(context.Background(), newImg)
	require.NoError(t, err)
	require.NotEmpty(t, newImgData)
	require.Equal(t, newImgData.PostID, imgData.PostID)
	require.NotEqual(t, newImgData.Page, imgData.Page)
	require.NotEqual(t, newImgData.Img1, imgData.Img1)
	require.NotEqual(t, newImgData.Img2, imgData.Img2)
	require.NotEqual(t, newImgData.Img3, imgData.Img3)
	require.NotEqual(t, newImgData.Img4, imgData.Img4)
	require.NotEqual(t, newImgData.Img5, imgData.Img5)
	require.NotEmpty(t, newImgData.UpdatedAt)
}
