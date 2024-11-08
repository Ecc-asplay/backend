-- name: CreateComments :one
INSERT INTO COMMENTS (
    COMMENT_ID,
    USER_ID,
    POST_ID,
    STATUS,
    IS_PUBLIC,
    COMMENTS,
    REACTION,
    IS_CENSORED
) VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) RETURNING *;

-- name: GetCommentsList :many
SELECT
    *
FROM
    COMMENTS
WHERE
    POST_ID = $1
ORDER BY
    COMMENT_ID DESC;

-- name: UpdateComments :one
UPDATE COMMENTS
SET
    STATUS = $2,
    IS_PUBLIC = $3,
    COMMENTS = $4,
    REACTION = $5,
    IS_CENSORED = $6,
    UPDATED_AT = NOW(
    )
WHERE
    COMMENT_ID = $1 RETURNING *;

-- name: DeleteComments :exec
DELETE FROM COMMENTS
WHERE
    COMMENT_ID = $1;