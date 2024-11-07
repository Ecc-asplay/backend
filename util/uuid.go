package util

import (
	"fmt"

	"github.com/google/uuid"
)

func CreateUUID() uuid.UUID {
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Create UUID Error", err)
	}
	return id
}
