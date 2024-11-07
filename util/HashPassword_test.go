package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRandomHash(t *testing.T) {
	pw := RandomPassword(20)

	hash, err := Hash(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestCheckRandomHash(t *testing.T) {
	pw := RandomPassword(20)

	hash, err := Hash(pw)
	checked, err := CheckPassword(pw, hash)
	require.NoError(t, err)
	require.True(t, checked)
}
