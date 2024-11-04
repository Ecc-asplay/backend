package util

import (
	"github.com/google/uuid"
)

const (
	AdminID   = uint32(1)
	UserID    = uint32(2)
	PostID    = uint32(3)
	CommentID = uint32(4)
)

func CreateUUID(types string) uuid.UUID {
	switch types {
	case "admin":
		adminUUID, _ := uuid.NewDCESecurity(uuid.Group, AdminID)
		return adminUUID
	case "user":
		userUUID, _ := uuid.NewDCESecurity(uuid.Group, UserID)
		return userUUID
	case "post":
		postUUID, _ := uuid.NewDCESecurity(uuid.Group, PostID)
		return postUUID
	case "comment":
		commentUUID, _ := uuid.NewDCESecurity(uuid.Group, CommentID)
		return commentUUID
	}
	return uuid.Nil
}
