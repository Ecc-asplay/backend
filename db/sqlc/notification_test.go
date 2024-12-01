package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func CreateRandomNotification(t *testing.T, user User) Notification {

	notificationData := CreateNotificationParams{
		UserID:  user.UserID,
		Content: gofakeit.Sentence(20),
		Icon:    gofakeit.ImagePng(20, 20),
	}

	Notifi, err := testQueries.CreateNotification(context.Background(), notificationData)
	require.NoError(t, err)
	require.NotEmpty(t, Notifi)
	require.Equal(t, notificationData.UserID, Notifi.UserID)
	require.Equal(t, notificationData.Content, Notifi.Content)
	require.Equal(t, notificationData.Icon, Notifi.Icon)
	require.False(t, Notifi.IsRead)
	require.NotEmpty(t, Notifi.CreatedAt)

	return Notifi
}

func TestCreateNotification(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomNotification(t, user)
}

func TestGetNotification(t *testing.T) {
	user := CreateRandomUser(t)
	for i := 0; i < 20; i++ {
		CreateRandomNotification(t, user)
	}

	allNotif, err := testQueries.GetNotification(context.Background(), user.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, allNotif)
	require.GreaterOrEqual(t, len(allNotif), 1)
}

func TestUpdateNotification(t *testing.T) {
	user := CreateRandomUser(t)
	for i := 0; i < 20; i++ {
		CreateRandomNotification(t, user)
	}

	newData, err := testQueries.UpdateNotification(context.Background(), user.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, newData)
	for _, v := range newData {
		require.True(t, v.IsRead)
	}
}
