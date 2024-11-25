// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateAdminUser(ctx context.Context, arg CreateAdminUserParams) (Adminuser, error)
	CreateBlock(ctx context.Context, arg CreateBlockParams) (Blockuser, error)
	CreateBookmarks(ctx context.Context, arg CreateBookmarksParams) (Bookmark, error)
	CreateComments(ctx context.Context, arg CreateCommentsParams) (Comment, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateSearchedRecord(ctx context.Context, arg CreateSearchedRecordParams) (Searchrecord, error)
	CreateTag(ctx context.Context, arg CreateTagParams) (Tag, error)
	CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAdminUser(ctx context.Context, email string) error
	DeleteBookmarks(ctx context.Context, arg DeleteBookmarksParams) error
	DeleteComments(ctx context.Context, commentID uuid.UUID) error
	DeletePost(ctx context.Context, arg DeletePostParams) error
	DeleteUser(ctx context.Context, arg DeleteUserParams) error
	GetAllBlockUsersList(ctx context.Context) ([]Blockuser, error)
	GetAllBookmarks(ctx context.Context, userID uuid.UUID) ([]Bookmark, error)
	GetBlockUserlist(ctx context.Context, userID uuid.UUID) ([]Blockuser, error)
	GetCommentsList(ctx context.Context, postID uuid.UUID) ([]Comment, error)
	GetLogin(ctx context.Context, email string) (GetLoginRow, error)
	GetNotification(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	GetPasswordToAdminLogin(ctx context.Context, email string) (string, error)
	GetPostOfKeywords(ctx context.Context, dollar_1 string) ([]Post, error)
	GetPostsList(ctx context.Context) ([]Post, error)
	GetSearchedRecordList(ctx context.Context) ([]Searchrecord, error)
	GetSession(ctx context.Context, id uuid.UUID) (Token, error)
	GetTag(ctx context.Context, dollar_1 string) ([]string, error)
	GetUserAllPosts(ctx context.Context, userID uuid.UUID) ([]Post, error)
	GetUserData(ctx context.Context, userID uuid.UUID) (User, error)
	ResetPassword(ctx context.Context, arg ResetPasswordParams) error
	UnBlockUser(ctx context.Context, arg UnBlockUserParams) (Blockuser, error)
	UpdateComments(ctx context.Context, arg UpdateCommentsParams) (Comment, error)
	UpdateDiseaseAndCondition(ctx context.Context, arg UpdateDiseaseAndConditionParams) error
	UpdateEmail(ctx context.Context, arg UpdateEmailParams) error
	UpdateIsPrivacy(ctx context.Context, arg UpdateIsPrivacyParams) error
	UpdateName(ctx context.Context, arg UpdateNameParams) (User, error)
	UpdateNotification(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	UpdatePosts(ctx context.Context, arg UpdatePostsParams) (Post, error)
}

var _ Querier = (*Queries)(nil)
