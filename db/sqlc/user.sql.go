// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
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
) RETURNING user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
`

type CreateUserParams struct {
	UserID       uuid.UUID   `json:"user_id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	Birth        pgtype.Date `json:"birth"`
	Gender       string      `json:"gender"`
	Disease      string      `json:"disease"`
	Condition    string      `json:"condition"`
	Hashpassword string      `json:"hashpassword"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UserID,
		arg.Username,
		arg.Email,
		arg.Birth,
		arg.Gender,
		arg.Disease,
		arg.Condition,
		arg.Hashpassword,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Birth,
		&i.Gender,
		&i.IsPrivacy,
		&i.Disease,
		&i.Condition,
		&i.Hashpassword,
		&i.Certification,
		&i.ResetPasswordAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM USERS
WHERE
    USER_ID = $1
    AND EMAIL = $2
`

type DeleteUserParams struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func (q *Queries) DeleteUser(ctx context.Context, arg DeleteUserParams) error {
	_, err := q.db.Exec(ctx, deleteUser, arg.UserID, arg.Email)
	return err
}

const getLogin = `-- name: GetLogin :one
SELECT
    USER_ID,
    HASHPASSWORD
FROM
    USERS
WHERE
    EMAIL = $1 LIMIT 1
`

type GetLoginRow struct {
	UserID       uuid.UUID `json:"user_id"`
	Hashpassword string    `json:"hashpassword"`
}

func (q *Queries) GetLogin(ctx context.Context, email string) (GetLoginRow, error) {
	row := q.db.QueryRow(ctx, getLogin, email)
	var i GetLoginRow
	err := row.Scan(&i.UserID, &i.Hashpassword)
	return i, err
}

const getUserData = `-- name: GetUserData :one
SELECT
    user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
FROM
    USERS
WHERE
    USER_ID = $1 LIMIT 1
`

func (q *Queries) GetUserData(ctx context.Context, userID uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserData, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Birth,
		&i.Gender,
		&i.IsPrivacy,
		&i.Disease,
		&i.Condition,
		&i.Hashpassword,
		&i.Certification,
		&i.ResetPasswordAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const resetPassword = `-- name: ResetPassword :exec
UPDATE USERS
SET
    HASHPASSWORD = COALESCE($2, HASHPASSWORD),
    RESET_PASSWORD_AT = COALESCE($3, RESET_PASSWORD_AT)
WHERE
    USER_ID = $1 RETURNING user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
`

type ResetPasswordParams struct {
	UserID          uuid.UUID        `json:"user_id"`
	Hashpassword    string           `json:"hashpassword"`
	ResetPasswordAt pgtype.Timestamp `json:"reset_password_at"`
}

func (q *Queries) ResetPassword(ctx context.Context, arg ResetPasswordParams) error {
	_, err := q.db.Exec(ctx, resetPassword, arg.UserID, arg.Hashpassword, arg.ResetPasswordAt)
	return err
}

const updateDiseaseAndCondition = `-- name: UpdateDiseaseAndCondition :exec
UPDATE USERS
SET
    DISEASE = COALESCE($2, DISEASE),
    CONDITION = COALESCE($3, CONDITION),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
`

type UpdateDiseaseAndConditionParams struct {
	UserID    uuid.UUID `json:"user_id"`
	Disease   string    `json:"disease"`
	Condition string    `json:"condition"`
}

func (q *Queries) UpdateDiseaseAndCondition(ctx context.Context, arg UpdateDiseaseAndConditionParams) error {
	_, err := q.db.Exec(ctx, updateDiseaseAndCondition, arg.UserID, arg.Disease, arg.Condition)
	return err
}

const updateEmail = `-- name: UpdateEmail :exec
UPDATE USERS
SET
    EMAIL = COALESCE($2, EMAIL),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
`

type UpdateEmailParams struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func (q *Queries) UpdateEmail(ctx context.Context, arg UpdateEmailParams) error {
	_, err := q.db.Exec(ctx, updateEmail, arg.UserID, arg.Email)
	return err
}

const updateIsPrivacy = `-- name: UpdateIsPrivacy :exec
UPDATE USERS
SET
    IS_PRIVACY = COALESCE($2, IS_PRIVACY),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
`

type UpdateIsPrivacyParams struct {
	UserID    uuid.UUID `json:"user_id"`
	IsPrivacy bool      `json:"is_privacy"`
}

func (q *Queries) UpdateIsPrivacy(ctx context.Context, arg UpdateIsPrivacyParams) error {
	_, err := q.db.Exec(ctx, updateIsPrivacy, arg.UserID, arg.IsPrivacy)
	return err
}

const updateName = `-- name: UpdateName :one
UPDATE USERS
SET
    USERNAME = COALESCE($2, USERNAME),
    UPDATED_AT = NOW()
WHERE
    USER_ID = $1 RETURNING user_id, username, email, birth, gender, is_privacy, disease, condition, hashpassword, certification, reset_password_at, created_at, updated_at
`

type UpdateNameParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

func (q *Queries) UpdateName(ctx context.Context, arg UpdateNameParams) (User, error) {
	row := q.db.QueryRow(ctx, updateName, arg.UserID, arg.Username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.Birth,
		&i.Gender,
		&i.IsPrivacy,
		&i.Disease,
		&i.Condition,
		&i.Hashpassword,
		&i.Certification,
		&i.ResetPasswordAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
