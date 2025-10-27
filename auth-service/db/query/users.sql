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
SET password_hash = $2, updated_at = now(), change_pass_at = now()
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: AddEmail :exec
UPDATE users
SET email_encrypted = $1,
    email_hash = $2
WHERE user_id = $3;

-- name: GetEmail :one
SELECT email_encrypted, email_hash
FROM users
WHERE user_id = $1;

-- name: RemoveEmail :exec
UPDATE users
SET email_encrypted = NULL,
    email_hash = NULL
WHERE user_id = $1;

-- name: CheckActiveEmail :one
SELECT email_hash
FROM users
WHERE email_hash = $1;

-- name: Toggle2Fa :exec
UPDATE users
SET is_2fa = $2, updated_at = now()
WHERE user_id = $1;

-- name: UpdateDisplayName :exec
UPDATE users
SET display_name = $2, updated_at = now()
WHERE user_id = $1;

-- name: UpdateAvatarURL :exec
UPDATE users
SET avatar_url = $2, updated_at = now()
WHERE user_id = $1;

-- name: GetProfile :one
SELECT avatar_url, display_name, email_encrypted, is_2fa, created_at, updated_at
FROM users
WHERE user_id = $1;