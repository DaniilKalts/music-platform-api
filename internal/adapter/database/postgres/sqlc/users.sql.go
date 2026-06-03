
package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, email, username, password_hash, salt, role, subscription_type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, email, username, role, subscription_type, created_at, updated_at
`

type CreateUserParams struct {
	ID               uuid.UUID
	Email            string
	Username         string
	PasswordHash     string
	Salt             string
	Role             UserRole
	SubscriptionType SubscriptionType
}

type CreateUserRow struct {
	ID               uuid.UUID
	Email            string
	Username         string
	Role             UserRole
	SubscriptionType SubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Username,
		arg.PasswordHash,
		arg.Salt,
		arg.Role,
		arg.SubscriptionType,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Role,
		&i.SubscriptionType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id,
       email,
       username,
       role,
       subscription_type,
       created_at,
       updated_at
FROM users
WHERE id = $1
`

type GetUserByIDRow struct {
	ID               uuid.UUID
	Email            string
	Username         string
	Role             UserRole
	SubscriptionType SubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (GetUserByIDRow, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Role,
		&i.SubscriptionType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserCredentialsByEmail = `-- name: GetUserCredentialsByEmail :one
SELECT id,
       email,
       password_hash,
       salt,
       role
FROM users
WHERE email = $1
`

type GetUserCredentialsByEmailRow struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Salt         string
	Role         UserRole
}

func (q *Queries) GetUserCredentialsByEmail(ctx context.Context, email string) (GetUserCredentialsByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserCredentialsByEmail, email)
	var i GetUserCredentialsByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.Salt,
		&i.Role,
	)
	return i, err
}

const updateUserProfile = `-- name: UpdateUserProfile :one
UPDATE users
SET email      = COALESCE($1, email),
    username   = COALESCE($2, username),
    updated_at = NOW()
WHERE id = $3
RETURNING id, email, username, role, subscription_type, created_at, updated_at
`

type UpdateUserProfileParams struct {
	Email    pgtype.Text
	Username pgtype.Text
	ID       uuid.UUID
}

type UpdateUserProfileRow struct {
	ID               uuid.UUID
	Email            string
	Username         string
	Role             UserRole
	SubscriptionType SubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (UpdateUserProfileRow, error) {
	row := q.db.QueryRow(ctx, updateUserProfile, arg.Email, arg.Username, arg.ID)
	var i UpdateUserProfileRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Role,
		&i.SubscriptionType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserSubscription = `-- name: UpdateUserSubscription :one
UPDATE users
SET subscription_type = $2,
    updated_at        = NOW()
WHERE id = $1
RETURNING id, email, username, role, subscription_type, created_at, updated_at
`

type UpdateUserSubscriptionParams struct {
	ID               uuid.UUID
	SubscriptionType SubscriptionType
}

type UpdateUserSubscriptionRow struct {
	ID               uuid.UUID
	Email            string
	Username         string
	Role             UserRole
	SubscriptionType SubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) UpdateUserSubscription(ctx context.Context, arg UpdateUserSubscriptionParams) (UpdateUserSubscriptionRow, error) {
	row := q.db.QueryRow(ctx, updateUserSubscription, arg.ID, arg.SubscriptionType)
	var i UpdateUserSubscriptionRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Role,
		&i.SubscriptionType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
