-- name: CreateToken :one
INSERT INTO TOKEN (
    ID,
    EMAIL,
    ACCESS_TOKEN,
    ROLES,
    STATUS,
    EXPIRES_AT
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetSession :one
SELECT
    *
FROM
    TOKEN
WHERE
    ID = $1 LIMIT 1;