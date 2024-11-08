-- name: CreateBlock :one
INSERT INTO BLOCKUSER (
    USER_ID,
    BLOCKUSER_ID,
    REASON,
    STATUS
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetAllBlockUsersList :many
SELECT
    *
FROM
    BLOCKUSER
ORDER BY
    BLOCK_AT DESC;

-- name: GetBlockUserlist :many
SELECT
    *
FROM
    BLOCKUSER
WHERE
    USER_ID = $1
ORDER BY
    BLOCK_AT DESC;

-- name: UnBlockUser :one
UPDATE BLOCKUSER
SET
    STATUS = $3,
    UNBLOCK_AT = NOW(
    )
WHERE
    USER_ID = $1
    AND BLOCKUSER_ID = $2 RETURNING *;