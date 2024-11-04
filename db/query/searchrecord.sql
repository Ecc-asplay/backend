-- name: CreateSearchedRecord :one
INSERT INTO SEARCHRECORD (
    SEARCH_CONTENT,
    IS_USER
) VALUES(
    $1,
    $2
) RETURNING *;

-- name: GetSearchedRecordList :many
SELECT
    *
FROM
    SEARCHRECORD
ORDER BY
    SEARCHED_AT DESC;

-- name: GetKeyWordSearchedRecord :many
SELECT
    *
FROM
    SEARCHRECORD
WHERE
    SEARCH_CONTENT LIKE '%'
                        || $1
                        || '%';