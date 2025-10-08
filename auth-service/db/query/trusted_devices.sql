-- name: CreateOrUpdateTrustedDevice :exec
INSERT INTO trusted_devices (user_id, device_id, browser, os, expires_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, device_id)
DO UPDATE SET
  expires_at = EXCLUDED.expires_at,
  updated_at = now();

-- name: GetTrustedDevices :many
SELECT id, device_id, browser, os, expires_at, created_at, updated_at
FROM trusted_devices
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: IsValidTrustedDevice :one
SELECT 1
FROM trusted_devices
WHERE user_id = $1 AND device_id = $2 AND browser = $3 AND os = $4 AND expires_at > now();

-- name: DeleteTrustedDevice :exec
DELETE FROM trusted_devices
WHERE user_id = $1 AND device_id = $2;
