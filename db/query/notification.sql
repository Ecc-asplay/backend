-- name: CreateNotification :one
INSERT INTO NOTIFICATION (
    USER_ID,
    CONTENT,
    ICON
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetNotification :many
SELECT
    *
FROM
    NOTIFICATION
WHERE
    USER_ID = $1;

-- name: UpdateNotification :many
UPDATE NOTIFICATION
SET
    IS_READ = TRUE
WHERE
    USER_ID = $1 RETURNING *;