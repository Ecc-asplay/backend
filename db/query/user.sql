-- name: CreateUser :one
INSERT INTO USERS (
    USERNAME,
    EMAIL,
    BIRTH,
    GENDER,
    DISEASE,
    CONDITION,
    HASHPASSWORD
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
) RETURNING *;

-- name: GetUserData :one
SELECT
    *
FROM
    USERS
WHERE
    USER_ID = $1 LIMIT 1;

-- name: GetPasswordToUserLogin :one
SELECT
    HASHPASSWORD
FROM
    USERS
WHERE
    EMAIL = $1 LIMIT 1;

-- name: UpdateName :one
UPDATE USERS
SET
    USERNAME = $2
WHERE
    USER_ID = $1 RETURNING *;

-- name: UpdateDiseaseAndCondition :exec
UPDATE USERS
SET
    DISEASE = $2,
    CONDITION = $3
WHERE
    USER_ID = $1 RETURNING *;

-- name: UpdateIsPrivacy :exec
UPDATE USERS
SET
    IS_PRIVACY = $2
WHERE
    USER_ID = $1 RETURNING *;

-- name: UpdateEmail :exec
UPDATE USERS
SET
    EMAIL = $2
WHERE
    USER_ID = $1 RETURNING *;

-- name: ResetPassword :exec
UPDATE USERS
SET
    HASHPASSWORD $2,
    RESET_PASSWORD_AT = $3
WHERE
    USER_ID = $1;

-- name: DeleteUser
DELETE FROM USERS
WHERE
    USER_ID = $1,
    EMAIL = $2;