-- name: CreateTag :one
INSERT INTO TAG (
    POST_ID,
    TAG_COMMENTS
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: FindTag :many
SELECT
    TAG_COMMENTS
FROM
    TAG
WHERE
    TAG_COMMENTS LIKE '%'
                      || CAST($1 AS TEXT)
                      || '%';