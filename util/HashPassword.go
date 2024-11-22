package util

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func Hash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		ErrorLog(err)
		return "", err
	}

	return hash, nil
}

func CheckPassword(input, hashedPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(input, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}

	return match, nil
}
