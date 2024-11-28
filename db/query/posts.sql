-- name: CreatePost :one
INSERT INTO POSTS (
    POST_ID,
    USER_ID,
    SHOW_ID,
    TITLE,
    FEEL,
    CONTENT,
    REACTION,
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
) RETURNING *;

-- name: GetUserAllPosts :many
SELECT
    *
FROM
    POSTS
WHERE
    USER_ID = $1;

-- name: GetPostsList :many
SELECT
    *
FROM
    POSTS
ORDER BY
    CREATED_AT DESC;

-- name: GetPostOfKeywords :many
SELECT
    *
FROM
    POSTS
WHERE
    TITLE LIKE '%'
               || CAST($1 AS TEXT)
               || '%'
    OR EXISTS (
        SELECT
            1
        FROM
            JSONB_ARRAY_ELEMENTS(CONTENT) AS ELEM,
            JSONB_ARRAY_ELEMENTS(ELEM->'children') AS CHILD
        WHERE
            CHILD->>'text' LIKE '%'
                                || CAST($1 AS TEXT)
                                || '%'
    );

-- name: UpdatePosts :one
UPDATE POSTS
SET
    SHOW_ID = COALESCE($3, SHOW_ID),
    TITLE = COALESCE($4, TITLE),
    FEEL = COALESCE($5, FEEL),
    CONTENT = COALESCE($6, CONTENT),
    REACTION = COALESCE($7, REACTION),
    IS_SENSITIVE = COALESCE($8, IS_SENSITIVE),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1
    AND POST_ID = $2
RETURNING *;

-- name: DeletePost :exec
DELETE FROM POSTS
WHERE
    USER_ID = $1
    AND POST_ID = $2 RETURNING *;