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
    STATUS = COALESCE($2, STATUS),
    IS_PUBLIC = COALESCE($3, IS_PUBLIC),
    COMMENTS = COALESCE($4, COMMENTS),
    REACTION = COALESCE($5, REACTION),
    IS_CENSORED = COALESCE($6, IS_CENSORED),
    UPDATED_AT = NOW()
WHERE
    COMMENT_ID = $1
RETURNING *;

-- name: DeleteComments :exec
DELETE FROM COMMENTS
WHERE
    COMMENT_ID = $1;