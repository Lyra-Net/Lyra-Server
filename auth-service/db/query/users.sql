-- name: CreateUser :exec
INSERT INTO users (user_id, display_name, username, password_hash)
VALUES ($1, $2, $3, $4);

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1;

-- name: UpdatePassword :exec
UPDATE users
SET password_hash = $2, updated_at = now()
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: AddEmail :exec
UPDATE users
SET email_encrypted = $1,
    email_hash = $2
WHERE user_id = $3;

-- name: CheckActiveEmail :one
SELECT email_hash
FROM users
WHERE email_hash = $1;