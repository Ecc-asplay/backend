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