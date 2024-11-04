// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bookmarks.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createBookmarks = `-- name: CreateBookmarks :one
INSERT INTO BOOKMARKS (
    POST_ID,
    USER_ID
) VALUES(
    $1,
    $2
) RETURNING user_id, post_id, created_at
`

type CreateBookmarksParams struct {
	PostID uuid.UUID `json:"post_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateBookmarks(ctx context.Context, arg CreateBookmarksParams) (Bookmark, error) {
	row := q.db.QueryRow(ctx, createBookmarks, arg.PostID, arg.UserID)
	var i Bookmark
	err := row.Scan(&i.UserID, &i.PostID, &i.CreatedAt)
	return i, err
}

const deleteBookmarks = `-- name: DeleteBookmarks :one
DELETE FROM BOOKMARKS
WHERE
    USER_ID = $1
    AND POST_ID = $2 RETURNING user_id, post_id, created_at
`

type DeleteBookmarksParams struct {
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
}

func (q *Queries) DeleteBookmarks(ctx context.Context, arg DeleteBookmarksParams) (Bookmark, error) {
	row := q.db.QueryRow(ctx, deleteBookmarks, arg.UserID, arg.PostID)
	var i Bookmark
	err := row.Scan(&i.UserID, &i.PostID, &i.CreatedAt)
	return i, err
}

const getBookmarks = `-- name: GetBookmarks :many
SELECT
    user_id, post_id, created_at
FROM
    BOOKMARKS
WHERE
    USER_ID = $1
ORDER BY
    CREATED_AT DESC
`

func (q *Queries) GetBookmarks(ctx context.Context, userID uuid.UUID) ([]Bookmark, error) {
	rows, err := q.db.Query(ctx, getBookmarks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Bookmark{}
	for rows.Next() {
		var i Bookmark
		if err := rows.Scan(&i.UserID, &i.PostID, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
