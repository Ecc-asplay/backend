-- name: CreateImage :one
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
) RETURNING *;

-- name: GetImage :many
SELECT
    *
FROM
    images
WHERE
    post_id = $1;

-- name: UpdateImage :one
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
    post_id = $1 RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images
WHERE
    post_id = $1;
