// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: comments_reaction.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createCommentsReaction = `-- name: CreateCommentsReaction :one
INSERT INTO COMMENTS_REACTION (
    USER_ID,
    COMMENT_ID,
    C_REACTION_THANKS,
    C_REACTION_HELPFUL,
    C_REACTION_USEFUL,
    C_REACTION_HEART
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at
`

type CreateCommentsReactionParams struct {
	UserID           uuid.UUID `json:"user_id"`
	CommentID        uuid.UUID `json:"comment_id"`
	CReactionThanks  bool      `json:"c_reaction_thanks"`
	CReactionHelpful bool      `json:"c_reaction_helpful"`
	CReactionUseful  bool      `json:"c_reaction_useful"`
	CReactionHeart   bool      `json:"c_reaction_heart"`
}

func (q *Queries) CreateCommentsReaction(ctx context.Context, arg CreateCommentsReactionParams) (CommentsReaction, error) {
	row := q.db.QueryRow(ctx, createCommentsReaction,
		arg.UserID,
		arg.CommentID,
		arg.CReactionThanks,
		arg.CReactionHelpful,
		arg.CReactionUseful,
		arg.CReactionHeart,
	)
	var i CommentsReaction
	err := row.Scan(
		&i.UserID,
		&i.CommentID,
		&i.CReactionThanks,
		&i.CReactionHeart,
		&i.CReactionHelpful,
		&i.CReactionUseful,
		&i.CreatedAt,
	)
	return i, err
}

const deleteCommentsReaction = `-- name: DeleteCommentsReaction :exec
DELETE FROM COMMENTS_REACTION
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
    AND NOT (C_REACTION_THANKS OR C_REACTION_HEART OR C_REACTION_HELPFUL OR C_REACTION_USEFUL)
`

type DeleteCommentsReactionParams struct {
	UserID    uuid.UUID `json:"user_id"`
	CommentID uuid.UUID `json:"comment_id"`
}

func (q *Queries) DeleteCommentsReaction(ctx context.Context, arg DeleteCommentsReactionParams) error {
	_, err := q.db.Exec(ctx, deleteCommentsReaction, arg.UserID, arg.CommentID)
	return err
}

const getAllCommentsReactionData = `-- name: GetAllCommentsReactionData :many
SELECT user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at FROM COMMENTS_REACTION
`

func (q *Queries) GetAllCommentsReactionData(ctx context.Context) ([]CommentsReaction, error) {
	rows, err := q.db.Query(ctx, getAllCommentsReactionData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CommentsReaction{}
	for rows.Next() {
		var i CommentsReaction
		if err := rows.Scan(
			&i.UserID,
			&i.CommentID,
			&i.CReactionThanks,
			&i.CReactionHeart,
			&i.CReactionHelpful,
			&i.CReactionUseful,
			&i.CreatedAt,
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

const getCommentsHeartOfTrue = `-- name: GetCommentsHeartOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_HEART = TRUE AND COMMENT_ID = $1
`

func (q *Queries) GetCommentsHeartOfTrue(ctx context.Context, commentID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getCommentsHeartOfTrue, commentID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCommentsHelpfulOfTrue = `-- name: GetCommentsHelpfulOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_HELPFUL = TRUE AND COMMENT_ID = $1
`

func (q *Queries) GetCommentsHelpfulOfTrue(ctx context.Context, commentID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getCommentsHelpfulOfTrue, commentID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCommentsReaction = `-- name: GetCommentsReaction :many
SELECT user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at FROM COMMENTS_REACTION 
WHERE COMMENT_ID = $1
`

func (q *Queries) GetCommentsReaction(ctx context.Context, commentID uuid.UUID) ([]CommentsReaction, error) {
	rows, err := q.db.Query(ctx, getCommentsReaction, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CommentsReaction{}
	for rows.Next() {
		var i CommentsReaction
		if err := rows.Scan(
			&i.UserID,
			&i.CommentID,
			&i.CReactionThanks,
			&i.CReactionHeart,
			&i.CReactionHelpful,
			&i.CReactionUseful,
			&i.CreatedAt,
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

const getCommentsThanksOfTrue = `-- name: GetCommentsThanksOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_THANKS = TRUE AND COMMENT_ID = $1
`

func (q *Queries) GetCommentsThanksOfTrue(ctx context.Context, commentID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getCommentsThanksOfTrue, commentID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCommentsUsefulOfTrue = `-- name: GetCommentsUsefulOfTrue :one
SELECT COUNT(*)
FROM COMMENTS_REACTION
WHERE C_REACTION_USEFUL = TRUE AND COMMENT_ID = $1
`

func (q *Queries) GetCommentsUsefulOfTrue(ctx context.Context, commentID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getCommentsUsefulOfTrue, commentID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateCommentsReactionHeart = `-- name: UpdateCommentsReactionHeart :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_HEART = NOT C_REACTION_HEART
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at
`

type UpdateCommentsReactionHeartParams struct {
	UserID    uuid.UUID `json:"user_id"`
	CommentID uuid.UUID `json:"comment_id"`
}

func (q *Queries) UpdateCommentsReactionHeart(ctx context.Context, arg UpdateCommentsReactionHeartParams) (CommentsReaction, error) {
	row := q.db.QueryRow(ctx, updateCommentsReactionHeart, arg.UserID, arg.CommentID)
	var i CommentsReaction
	err := row.Scan(
		&i.UserID,
		&i.CommentID,
		&i.CReactionThanks,
		&i.CReactionHeart,
		&i.CReactionHelpful,
		&i.CReactionUseful,
		&i.CreatedAt,
	)
	return i, err
}

const updateCommentsReactionHelpful = `-- name: UpdateCommentsReactionHelpful :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_HELPFUL = NOT C_REACTION_HELPFUL
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at
`

type UpdateCommentsReactionHelpfulParams struct {
	UserID    uuid.UUID `json:"user_id"`
	CommentID uuid.UUID `json:"comment_id"`
}

func (q *Queries) UpdateCommentsReactionHelpful(ctx context.Context, arg UpdateCommentsReactionHelpfulParams) (CommentsReaction, error) {
	row := q.db.QueryRow(ctx, updateCommentsReactionHelpful, arg.UserID, arg.CommentID)
	var i CommentsReaction
	err := row.Scan(
		&i.UserID,
		&i.CommentID,
		&i.CReactionThanks,
		&i.CReactionHeart,
		&i.CReactionHelpful,
		&i.CReactionUseful,
		&i.CreatedAt,
	)
	return i, err
}

const updateCommentsReactionThanks = `-- name: UpdateCommentsReactionThanks :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_THANKS = NOT C_REACTION_THANKS
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at
`

type UpdateCommentsReactionThanksParams struct {
	UserID    uuid.UUID `json:"user_id"`
	CommentID uuid.UUID `json:"comment_id"`
}

func (q *Queries) UpdateCommentsReactionThanks(ctx context.Context, arg UpdateCommentsReactionThanksParams) (CommentsReaction, error) {
	row := q.db.QueryRow(ctx, updateCommentsReactionThanks, arg.UserID, arg.CommentID)
	var i CommentsReaction
	err := row.Scan(
		&i.UserID,
		&i.CommentID,
		&i.CReactionThanks,
		&i.CReactionHeart,
		&i.CReactionHelpful,
		&i.CReactionUseful,
		&i.CreatedAt,
	)
	return i, err
}

const updateCommentsReactionUseful = `-- name: UpdateCommentsReactionUseful :one
UPDATE COMMENTS_REACTION
SET 
    C_REACTION_USEFUL = NOT C_REACTION_USEFUL
WHERE 
    USER_ID = $1 AND COMMENT_ID = $2
RETURNING user_id, comment_id, c_reaction_thanks, c_reaction_heart, c_reaction_helpful, c_reaction_useful, created_at
`

type UpdateCommentsReactionUsefulParams struct {
	UserID    uuid.UUID `json:"user_id"`
	CommentID uuid.UUID `json:"comment_id"`
}

func (q *Queries) UpdateCommentsReactionUseful(ctx context.Context, arg UpdateCommentsReactionUsefulParams) (CommentsReaction, error) {
	row := q.db.QueryRow(ctx, updateCommentsReactionUseful, arg.UserID, arg.CommentID)
	var i CommentsReaction
	err := row.Scan(
		&i.UserID,
		&i.CommentID,
		&i.CReactionThanks,
		&i.CReactionHeart,
		&i.CReactionHelpful,
		&i.CReactionUseful,
		&i.CreatedAt,
	)
	return i, err
}
