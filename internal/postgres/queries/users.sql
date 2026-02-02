-- name: CreateUser :one
INSERT INTO users (id, email, password_hash,role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: UpdateUser :one
UPDATE users
SET email = $2, password_hash = $3, role = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;