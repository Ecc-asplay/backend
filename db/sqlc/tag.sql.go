// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: tag.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createTag = `-- name: CreateTag :one
INSERT INTO TAG (
    POST_ID,
    TAG_COMMENTS
) VALUES (
    $1,
    $2
) RETURNING post_id, tag_comments
`

type CreateTagParams struct {
	PostID      uuid.UUID `json:"post_id"`
	TagComments string    `json:"tag_comments"`
}

func (q *Queries) CreateTag(ctx context.Context, arg CreateTagParams) (Tag, error) {
	row := q.db.QueryRow(ctx, createTag, arg.PostID, arg.TagComments)
	var i Tag
	err := row.Scan(&i.PostID, &i.TagComments)
	return i, err
}

const findTag = `-- name: FindTag :many
SELECT
    TAG_COMMENTS
FROM
    TAG
WHERE
    TAG_COMMENTS LIKE '%'
                      || CAST($1 AS TEXT)
                      || '%'
`

func (q *Queries) FindTag(ctx context.Context, dollar_1 string) ([]string, error) {
	rows, err := q.db.Query(ctx, findTag, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var tag_comments string
		if err := rows.Scan(&tag_comments); err != nil {
			return nil, err
		}
		items = append(items, tag_comments)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
