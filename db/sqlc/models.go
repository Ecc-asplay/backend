// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Adminuser struct {
	Email        string           `json:"email"`
	Hashpassword string           `json:"hashpassword"`
	StaffName    string           `json:"staff_name"`
	Department   string           `json:"department"`
	JoinedAt     pgtype.Timestamp `json:"joined_at"`
}

type Blockuser struct {
	UserID      uuid.UUID        `json:"user_id"`
	BlockuserID uuid.UUID        `json:"blockuser_id"`
	Reason      string           `json:"reason"`
	Status      string           `json:"status"`
	BlockAt     pgtype.Timestamp `json:"block_at"`
	UnblockAt   pgtype.Timestamp `json:"unblock_at"`
}

type Bookmark struct {
	UserID    uuid.UUID        `json:"user_id"`
	PostID    uuid.UUID        `json:"post_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Comment struct {
	CommentID  uuid.UUID        `json:"comment_id"`
	UserID     uuid.UUID        `json:"user_id"`
	PostID     uuid.UUID        `json:"post_id"`
	Status     string           `json:"status"`
	IsPublic   bool             `json:"is_public"`
	Comments   string           `json:"comments"`
	Reaction   int32            `json:"reaction"`
	IsCensored bool             `json:"is_censored"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}

type Notification struct {
	UserID    uuid.UUID        `json:"user_id"`
	Content   string           `json:"content"`
	Icon      []byte           `json:"icon"`
	IsRead    bool             `json:"is_read"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Post struct {
	UserID      uuid.UUID        `json:"user_id"`
	PostID      uuid.UUID        `json:"post_id"`
	ShowID      string           `json:"show_id"`
	Title       string           `json:"title"`
	Feel        string           `json:"feel"`
	Content     []byte           `json:"content"`
	Reaction    int32            `json:"reaction"`
	IsSensitive bool             `json:"is_sensitive"`
	Status      string           `json:"status"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type Searchrecord struct {
	SearchContent string           `json:"search_content"`
	IsUser        bool             `json:"is_user"`
	SearchedAt    pgtype.Timestamp `json:"searched_at"`
}

type Tag struct {
	PostID      uuid.UUID `json:"post_id"`
	TagComments string    `json:"tag_comments"`
}

type Token struct {
	ID          uuid.UUID        `json:"id"`
	UserID      uuid.UUID        `json:"user_id"`
	AccessToken string           `json:"access_token"`
	Roles       string           `json:"roles"`
	Status      string           `json:"status"`
	TakeAt      pgtype.Timestamp `json:"take_at"`
	ExpiresAt   pgtype.Timestamp `json:"expires_at"`
}

type User struct {
	UserID          uuid.UUID        `json:"user_id"`
	Username        string           `json:"username"`
	Email           string           `json:"email"`
	Birth           pgtype.Date      `json:"birth"`
	Gender          string           `json:"gender"`
	IsPrivacy       bool             `json:"is_privacy"`
	Disease         string           `json:"disease"`
	Condition       string           `json:"condition"`
	Hashpassword    string           `json:"hashpassword"`
	Certification   bool             `json:"certification"`
	ResetPasswordAt pgtype.Timestamp `json:"reset_password_at"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
}
