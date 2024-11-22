package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomHash(t *testing.T, pw string) string {
	hash, err := Hash(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	return hash
}

func TestCreateRandomHash(t *testing.T) {
	pw := RandomPassword(20)
	CreateRandomHash(t, pw)
}

func TestCheckRandomHash(t *testing.T) {
	pw := RandomPassword(20)
	hash := CreateRandomHash(t, pw)

	checked, err := CheckPassword(pw, hash)
	require.NoError(t, err)
	require.True(t, checked)
}
