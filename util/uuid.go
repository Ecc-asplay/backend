package util

import (
	"github.com/google/uuid"
)

func CreateUUID() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		ErrorLog(err)
	}
	return id
}
