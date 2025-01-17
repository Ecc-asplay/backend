// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: comments.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createComments = `-- name: CreateComments :one
INSERT INTO COMMENTS (
    COMMENT_ID,
    USER_ID,
    POST_ID,
    POST_USER,
    STATUS,
    IS_PUBLIC,
    COMMENTS,
    IS_CENSORED
) VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) RETURNING comment_id, user_id, post_id, status, is_public, comments, is_censored, created_at, updated_at, post_user
`

type CreateCommentsParams struct {
	CommentID  uuid.UUID `json:"comment_id"`
	UserID     uuid.UUID `json:"user_id"`
	PostID     uuid.UUID `json:"post_id"`
	PostUser   uuid.UUID `json:"post_user"`
	Status     string    `json:"status"`
	IsPublic   bool      `json:"is_public"`
	Comments   string    `json:"comments"`
	IsCensored bool      `json:"is_censored"`
}

func (q *Queries) CreateComments(ctx context.Context, arg CreateCommentsParams) (Comment, error) {
	row := q.db.QueryRow(ctx, createComments,
		arg.CommentID,
		arg.UserID,
		arg.PostID,
		arg.PostUser,
		arg.Status,
		arg.IsPublic,
		arg.Comments,
		arg.IsCensored,
	)
	var i Comment
	err := row.Scan(
		&i.CommentID,
		&i.UserID,
		&i.PostID,
		&i.Status,
		&i.IsPublic,
		&i.Comments,
		&i.IsCensored,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PostUser,
	)
	return i, err
}

const deleteComments = `-- name: DeleteComments :exec
DELETE FROM COMMENTS
WHERE
    COMMENT_ID = $1
`

func (q *Queries) DeleteComments(ctx context.Context, commentID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteComments, commentID)
	return err
}

const getAllComments = `-- name: GetAllComments :many
SELECT
    comment_id, user_id, post_id, status, is_public, comments, is_censored, created_at, updated_at, post_user
FROM
    COMMENTS
WHERE
    post_user = $1
ORDER BY
    COMMENT_ID DESC
`

func (q *Queries) GetAllComments(ctx context.Context, postUser uuid.UUID) ([]Comment, error) {
	rows, err := q.db.Query(ctx, getAllComments, postUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Comment{}
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.CommentID,
			&i.UserID,
			&i.PostID,
			&i.Status,
			&i.IsPublic,
			&i.Comments,
			&i.IsCensored,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PostUser,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCommentsList = `-- name: GetCommentsList :many
SELECT
    comment_id, user_id, post_id, status, is_public, comments, is_censored, created_at, updated_at, post_user
FROM
    COMMENTS
WHERE
    POST_ID = $1
ORDER BY
    COMMENT_ID DESC
`

func (q *Queries) GetCommentsList(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	rows, err := q.db.Query(ctx, getCommentsList, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Comment{}
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.CommentID,
			&i.UserID,
			&i.PostID,
			&i.Status,
			&i.IsPublic,
			&i.Comments,
			&i.IsCensored,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PostUser,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateComments = `-- name: UpdateComments :one
UPDATE COMMENTS
SET
    STATUS = COALESCE($2, STATUS),
    IS_PUBLIC = COALESCE($3, IS_PUBLIC),
    COMMENTS = COALESCE($4, COMMENTS),
    IS_CENSORED = COALESCE($5, IS_CENSORED),
    UPDATED_AT = NOW()
WHERE
    COMMENT_ID = $1
RETURNING comment_id, user_id, post_id, status, is_public, comments, is_censored, created_at, updated_at, post_user
`

type UpdateCommentsParams struct {
	CommentID  uuid.UUID `json:"comment_id"`
	Status     string    `json:"status"`
	IsPublic   bool      `json:"is_public"`
	Comments   string    `json:"comments"`
	IsCensored bool      `json:"is_censored"`
}

func (q *Queries) UpdateComments(ctx context.Context, arg UpdateCommentsParams) (Comment, error) {
	row := q.db.QueryRow(ctx, updateComments,
		arg.CommentID,
		arg.Status,
		arg.IsPublic,
		arg.Comments,
		arg.IsCensored,
	)
	var i Comment
	err := row.Scan(
		&i.CommentID,
		&i.UserID,
		&i.PostID,
		&i.Status,
		&i.IsPublic,
		&i.Comments,
		&i.IsCensored,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PostUser,
	)
	return i, err
}
