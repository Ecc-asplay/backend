<<<<<<< HEAD
-- name: CreateToken :one
INSERT INTO TOKEN (
    ID,
    user_id,
    ACCESS_TOKEN,
    ROLES,
    STATUS,
    EXPIRES_AT
) VALUES (
    $1,    $2,
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
=======
-- name: CreateToken :one
INSERT INTO TOKEN (
    ID,
    USER_ID,
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
>>>>>>> main
    ID = $1 LIMIT 1;