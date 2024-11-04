-- name: CreatePost :exec
INSERT INTO POSTS (
    USER_ID,
    SHOW_ID,
    TITLE,
    FEEL,
    CONTENT,
    REACTION,
    IMAGE,
    IS_SENSITIVE,
    STATUS
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
);

-- name: GetPostForUser :one
SELECT
    *
FROM
    POSTS
WHERE
    USER_ID = $1 RETURNING *;

-- name: GetPostsList :many
SELECT
    *
FROM
    POSTS
ORDER BY
    CREATED_AT DESC RETURNING *;

-- name: UpdatePosts :exec
UPDATE POSTS
SET
    SHOW_ID = $3,
    TITLE = $4,
    FEEL = $5,
    CONTENT = $6,
    REACTION = $7,
    IMAGE = $8,
    IS_SENSITIVE = $9
WHERE
    USER_ID = $1
    AND POST_ID = $2 RETURNING *;

-- name: DeletePost :exec
DELETE FROM POSTS
WHERE
    USER_ID = $1
    AND POST_ID = $2;