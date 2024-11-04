--name: CreateBookmark :one
INSERT INTO BOOKMARK (
    POST_ID,
    USER_ID
) VALUES(
    $1,
    $2
) RETURNING *;

--name: GetBookmarks :many
SELECT
    *
FROM
    BOOKMARK
WHERE
    USER_ID = $1
ORDER BY
    CREATED_AT DESC RETURNING *;

--name: DeleteBookmarks :one
DELETE FROM BOOKMARK
WHERE
    USER_ID = $1
    AND POST_ID = $2 RETURNING *;