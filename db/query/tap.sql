--name: CreateTag :one
INSERT INTO TAG (
    POST_ID,
    TAG
) VALUES (
    $1,
    $2
) RETURNING *;

--name: GetTag :many
SELECT
    *
FROM
    TAG
WHERE
    TAG LIKE '%'
             || $1
             || '%';