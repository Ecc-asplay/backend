package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Ecc-asplay/backend/util"
)

func TestCreateRandomPaseto(t *testing.T) {
	symmetricKey := util.RandomPassword(32)
	Maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotEmpty(t, Maker)

	userid := util.CreateUUID()
	role := util.RandomRole()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, payload, err := Maker.CreateToken(userid, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.GreaterOrEqual(t, len(token), 32)
	require.NotEmpty(t, payload.ID)
	require.GreaterOrEqual(t, len(payload.ID), 16)
	require.NotEmpty(t, payload.IssuedAt)
	require.Equal(t, payload.UserID, userid)
	require.Equal(t, payload.Role, role)

	payload, err = Maker.VerifyToken(token)
	require.NoError(t, err)

	require.NotZero(t, payload.ID)
	require.Equal(t, userid, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
