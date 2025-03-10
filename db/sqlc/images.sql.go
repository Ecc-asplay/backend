// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: images.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images (
    post_id,
    page,
    img1,
    img2,
    img3,
    img4,
    img5
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING post_id, page, img1, img2, img3, img4, img5, created_at, updated_at
`

type CreateImageParams struct {
	PostID uuid.UUID `json:"post_id"`
	Page   int32     `json:"page"`
	Img1   []byte    `json:"img1"`
	Img2   []byte    `json:"img2"`
	Img3   []byte    `json:"img3"`
	Img4   []byte    `json:"img4"`
	Img5   []byte    `json:"img5"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, createImage,
		arg.PostID,
		arg.Page,
		arg.Img1,
		arg.Img2,
		arg.Img3,
		arg.Img4,
		arg.Img5,
	)
	var i Image
	err := row.Scan(
		&i.PostID,
		&i.Page,
		&i.Img1,
		&i.Img2,
		&i.Img3,
		&i.Img4,
		&i.Img5,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteImage = `-- name: DeleteImage :exec
DELETE FROM images
WHERE
    post_id = $1
`

func (q *Queries) DeleteImage(ctx context.Context, postID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteImage, postID)
	return err
}

const getImage = `-- name: GetImage :many
SELECT
    post_id, page, img1, img2, img3, img4, img5, created_at, updated_at
FROM
    images
WHERE
    post_id = $1
`

func (q *Queries) GetImage(ctx context.Context, postID uuid.UUID) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImage, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.PostID,
			&i.Page,
			&i.Img1,
			&i.Img2,
			&i.Img3,
			&i.Img4,
			&i.Img5,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateImage = `-- name: UpdateImage :one
UPDATE images
SET
    page = COALESCE($2, page),
    img1 = COALESCE($3, img1),
    img2 = COALESCE($4, img2),
    img3 = COALESCE($5, img3),
    img4 = COALESCE($6, img4),
    img5 = COALESCE($7, img5),
    updated_at = Now()
WHERE
    post_id = $1 RETURNING post_id, page, img1, img2, img3, img4, img5, created_at, updated_at
`

type UpdateImageParams struct {
	PostID uuid.UUID `json:"post_id"`
	Page   int32     `json:"page"`
	Img1   []byte    `json:"img1"`
	Img2   []byte    `json:"img2"`
	Img3   []byte    `json:"img3"`
	Img4   []byte    `json:"img4"`
	Img5   []byte    `json:"img5"`
}

func (q *Queries) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, updateImage,
		arg.PostID,
		arg.Page,
		arg.Img1,
		arg.Img2,
		arg.Img3,
		arg.Img4,
		arg.Img5,
	)
	var i Image
	err := row.Scan(
		&i.PostID,
		&i.Page,
		&i.Img1,
		&i.Img2,
		&i.Img3,
		&i.Img4,
		&i.Img5,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
