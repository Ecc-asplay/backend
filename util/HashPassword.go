package util

import (
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
)

func Hash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, nil
}

func CheckPassword(hashedPassword, input string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(input, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}

	return match, nil
}
