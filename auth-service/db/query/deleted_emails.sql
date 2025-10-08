-- name: GetDeletedEmail :one
SELECT * FROM deleted_emails WHERE email_hash = $1;

