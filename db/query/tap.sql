-- name: CreateTag :one
INSERT INTO TAG (
    POST_ID,
    TAG_COMMENTS
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetTag :many
SELECT
    *
FROM
    TAG
WHERE
    TAG_COMMENTS LIKE '%'
                      || $1
                      || '%';