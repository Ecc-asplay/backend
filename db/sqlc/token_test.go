package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
)

func CreateRandomToken(t *testing.T) Token {
	make, err := token.NewPasetoMaker(util.RandomString(32))
	userid := util.CreateUUID()
	role := util.RandomRole()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := make.CreateToken(userid, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.ID)
	require.Equal(t, userid, payload.UserID)
	require.Equal(t, role, payload.Role)
	require.NotEqual(t, issuedAt, payload.IssuedAt)
	require.NotEqual(t, expiredAt, payload.ExpiredAt)

	tokenData := CreateTokenParams{
		ID:          payload.ID,
		UserID:      payload.UserID,
		AccessToken: token,
		Roles:       payload.Role,
		Status:      "OK",
		ExpiresAt: pgtype.Timestamp{
			Time:  payload.ExpiredAt,
			Valid: true,
		},
	}

	TokenDB, err := testQueries.CreateToken(context.Background(), tokenData)
	require.NoError(t, err)
	require.NotEmpty(t, TokenDB)
	require.Equal(t, tokenData.ID, TokenDB.ID)
	require.Equal(t, tokenData.UserID, TokenDB.UserID)
	require.Equal(t, tokenData.AccessToken, TokenDB.AccessToken)
	require.Equal(t, tokenData.Roles, TokenDB.Roles)
	require.Equal(t, tokenData.Status, TokenDB.Status)
	require.NotEqual(t, tokenData.ExpiresAt.Time, TokenDB.ExpiresAt.Time)
	require.NotEmpty(t, TokenDB.TakeAt)

	return TokenDB
}

func TestCreateRandomToken(t *testing.T) {
	CreateRandomToken(t)
}

func TestGetSession(t *testing.T) {
	token := CreateRandomToken(t)

	session, err := testQueries.GetSession(context.Background(), token.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, session.AccessToken, token.AccessToken)
	require.Equal(t, session.UserID, token.UserID)
	require.Equal(t, session.Roles, token.Roles)
	require.Equal(t, session.Status, token.Status)
	require.Equal(t, session.ExpiresAt.Time, token.ExpiresAt.Time)
	require.Equal(t, session.TakeAt.Time, token.TakeAt.Time)
}
