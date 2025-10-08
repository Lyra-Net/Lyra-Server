-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token, access_jti, device_id, browser, os, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE id = $1 AND user_id = $2;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;
