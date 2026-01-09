-- name: GetUserByTelegramID :one
SELECT * FROM users WHERE telegram_id = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (telegram_id, username, first_name, last_name, language_code)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = ?, first_name = ?, last_name = ?, language_code = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions WHERE token = ? LIMIT 1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < datetime('now');
