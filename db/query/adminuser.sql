-- name: CreateAdminUser :one
INSERT INTO ADMINUSER (
    admin_id,
    EMAIL,
    HASHPASSWORD,
    STAFF_NAME,
    DEPARTMENT
) VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetAdminLogin :one
SELECT
    admin_id,
    HASHPASSWORD
FROM
    ADMINUSER
WHERE
    EMAIL = $1 LIMIT 1;

-- name: DeleteAdminUser :exec
DELETE FROM ADMINUSER
WHERE
    EMAIL = $1;