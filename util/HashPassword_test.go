package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	pw := RandomPassword(20)

	hash, err := Hash(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	checked, err := CheckPassword(pw, hash)
	require.NoError(t, err)
	require.True(t, checked)
}
