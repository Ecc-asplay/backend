-- name: CreateUser :one
INSERT INTO USERS (
    USER_ID,
    USERNAME,
    EMAIL,
    BIRTH,
    GENDER,
    IS_PRIVACY,
    DISEASE,
    CONDITION,
    HASHPASSWORD,
    CERTIFICATION,
    RESET_PASSWORD_AT,
    CREATED_AT,
    UPDATED_AT
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
) RETURNING *;