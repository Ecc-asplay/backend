-- name: CreateBookmarks :one
INSERT INTO BOOKMARKS (
    POST_ID,
    USER_ID
) VALUES(
    $1,
    $2
) RETURNING *;

-- name: GetAllBookmarks :many
SELECT
    *
FROM
    BOOKMARKS
WHERE
    USER_ID = $1
ORDER BY
    CREATED_AT DESC;

-- name: DeleteBookmarks :exec
DELETE FROM BOOKMARKS
WHERE
    USER_ID = $1
    AND POST_ID = $2 RETURNING *;