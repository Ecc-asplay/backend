// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: blockuser.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createBlock = `-- name: CreateBlock :one
INSERT INTO BLOCKUSER (
    USER_ID,
    BLOCKUSER_ID,
    REASON,
    STATUS
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING user_id, blockuser_id, reason, status, block_at, unblock_at
`

type CreateBlockParams struct {
	UserID      uuid.UUID `json:"user_id"`
	BlockuserID uuid.UUID `json:"blockuser_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
}

func (q *Queries) CreateBlock(ctx context.Context, arg CreateBlockParams) (Blockuser, error) {
	row := q.db.QueryRow(ctx, createBlock,
		arg.UserID,
		arg.BlockuserID,
		arg.Reason,
		arg.Status,
	)
	var i Blockuser
	err := row.Scan(
		&i.UserID,
		&i.BlockuserID,
		&i.Reason,
		&i.Status,
		&i.BlockAt,
		&i.UnblockAt,
	)
	return i, err
}

const getAllBlockUsersList = `-- name: GetAllBlockUsersList :many
SELECT
    user_id, blockuser_id, reason, status, block_at, unblock_at
FROM
    BLOCKUSER
ORDER BY
    BLOCK_AT DESC
`

func (q *Queries) GetAllBlockUsersList(ctx context.Context) ([]Blockuser, error) {
	rows, err := q.db.Query(ctx, getAllBlockUsersList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Blockuser{}
	for rows.Next() {
		var i Blockuser
		if err := rows.Scan(
			&i.UserID,
			&i.BlockuserID,
			&i.Reason,
			&i.Status,
			&i.BlockAt,
			&i.UnblockAt,
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

const getBlockUserlist = `-- name: GetBlockUserlist :many
SELECT
    user_id, blockuser_id, reason, status, block_at, unblock_at
FROM
    BLOCKUSER
WHERE
    USER_ID = $1
ORDER BY
    BLOCK_AT DESC
`

func (q *Queries) GetBlockUserlist(ctx context.Context, userID uuid.UUID) ([]Blockuser, error) {
	rows, err := q.db.Query(ctx, getBlockUserlist, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Blockuser{}
	for rows.Next() {
		var i Blockuser
		if err := rows.Scan(
			&i.UserID,
			&i.BlockuserID,
			&i.Reason,
			&i.Status,
			&i.BlockAt,
			&i.UnblockAt,
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

const unBlockUser = `-- name: UnBlockUser :one
UPDATE BLOCKUSER
SET
    STATUS = COALESCE($3, STATUS),
    UNBLOCK_AT = NOW(
    )
WHERE
    USER_ID = $1
    AND BLOCKUSER_ID = $2 RETURNING user_id, blockuser_id, reason, status, block_at, unblock_at
`

type UnBlockUserParams struct {
	UserID      uuid.UUID `json:"user_id"`
	BlockuserID uuid.UUID `json:"blockuser_id"`
	Status      string    `json:"status"`
}

func (q *Queries) UnBlockUser(ctx context.Context, arg UnBlockUserParams) (Blockuser, error) {
	row := q.db.QueryRow(ctx, unBlockUser, arg.UserID, arg.BlockuserID, arg.Status)
	var i Blockuser
	err := row.Scan(
		&i.UserID,
		&i.BlockuserID,
		&i.Reason,
		&i.Status,
		&i.BlockAt,
		&i.UnblockAt,
	)
	return i, err
}
