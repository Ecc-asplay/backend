-- name: CreateAdminUser :one
INSERT INTO ADMINUSER (
    EMAIL,
    HASHPASSWORD,
    STAFF_NAME,
    DEPARTMENT
) VALUES(
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetPasswordToAdminLogin :one
SELECT
    HASHPASSWORD
FROM
    ADMINUSER
WHERE
    EMAIL = $1 LIMIT 1;

-- name: DeleteAdminUser :exec
DELETE FROM ADMINUSER
WHERE
    EMAIL = $1;