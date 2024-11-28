-- name: CreateUser :one
INSERT INTO USERS (
    USER_ID,
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
    $7,
    $8
) RETURNING *;

-- name: GetUserData :one
SELECT
    *
FROM
    USERS
WHERE
    USER_ID = $1 LIMIT 1;

-- name: GetLogin :one
SELECT
    USER_ID,
    HASHPASSWORD
FROM
    USERS
WHERE
    EMAIL = $1 LIMIT 1;

-- name: UpdateName :one
UPDATE USERS
SET
    USERNAME = COALESCE($2, USERNAME),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING *;

-- name: UpdateDiseaseAndCondition :exec
UPDATE USERS
SET
    DISEASE = COALESCE($2, DISEASE),
    CONDITION = COALESCE($3, CONDITION),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING *;

-- name: UpdateIsPrivacy :exec
UPDATE USERS
SET
    IS_PRIVACY = COALESCE($2, IS_PRIVACY),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING *;

-- name: UpdateEmail :exec
UPDATE USERS
SET
    EMAIL = COALESCE($2, EMAIL),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING *;

-- name: ResetPassword :exec
UPDATE USERS
SET
    HASHPASSWORD = COALESCE($2, HASHPASSWORD),
    RESET_PASSWORD_AT = COALESCE($3, RESET_PASSWORD_AT)
WHERE
    USER_ID = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM USERS
WHERE
    USER_ID = $1
    AND EMAIL = $2;