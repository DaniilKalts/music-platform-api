-- name: CreateUser :one
INSERT INTO users (id, email, username, password_hash, salt, role, subscription_type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, email, username, role, subscription_type, created_at, updated_at;

-- name: GetUserCredentialsByEmail :one
SELECT id,
       email,
       password_hash,
       salt,
       role
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id,
       email,
       username,
       role,
       subscription_type,
       created_at,
       updated_at
FROM users
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET email      = COALESCE(sqlc.narg('email'), email),
    username   = COALESCE(sqlc.narg('username'), username),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, email, username, role, subscription_type, created_at, updated_at;

-- name: UpdateUserSubscription :one
UPDATE users
SET subscription_type = $2,
    updated_at        = NOW()
WHERE id = $1
RETURNING id, email, username, role, subscription_type, created_at, updated_at;
