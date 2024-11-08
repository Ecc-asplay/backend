package util

import (
	"fmt"

	"github.com/google/uuid"
)

func CreateUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Create UUID Error", err)
	}
	return id
}
