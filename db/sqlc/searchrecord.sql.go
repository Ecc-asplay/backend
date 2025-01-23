// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: searchrecord.sql

package db

import (
	"context"
)

const createSearchedRecord = `-- name: CreateSearchedRecord :one
INSERT INTO SEARCHRECORD (
    SEARCH_CONTENT,
    IS_USER
) VALUES(
    $1,
    $2
) RETURNING search_content, is_user, searched_at
`

type CreateSearchedRecordParams struct {
	SearchContent string `json:"search_content"`
	IsUser        bool   `json:"is_user"`
}

func (q *Queries) CreateSearchedRecord(ctx context.Context, arg CreateSearchedRecordParams) (Searchrecord, error) {
	row := q.db.QueryRow(ctx, createSearchedRecord, arg.SearchContent, arg.IsUser)
	var i Searchrecord
	err := row.Scan(&i.SearchContent, &i.IsUser, &i.SearchedAt)
	return i, err
}

const getSearchedRecordList = `-- name: GetSearchedRecordList :many
SELECT
    search_content, is_user, searched_at
FROM
    SEARCHRECORD
ORDER BY
    SEARCHED_AT DESC
`

func (q *Queries) GetSearchedRecordList(ctx context.Context) ([]Searchrecord, error) {
	rows, err := q.db.Query(ctx, getSearchedRecordList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Searchrecord{}
	for rows.Next() {
		var i Searchrecord
		if err := rows.Scan(&i.SearchContent, &i.IsUser, &i.SearchedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
