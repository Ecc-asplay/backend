package util

import (
	"github.com/google/uuid"
)

func CreateUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		ErrorLog(err)
	}
	return id
}
