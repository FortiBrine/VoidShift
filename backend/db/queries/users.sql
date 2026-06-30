-- name: CreateUser :exec
INSERT INTO users (username, password_hash, admin)
VALUES (?, ?, ?)
ON CONFLICT(username) DO NOTHING;

-- name: GetUserByID :one
SELECT id, username, password_hash, admin FROM users WHERE id = ?;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, admin FROM users WHERE username = ?;
