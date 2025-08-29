-- name: CreateArtist :one
INSERT INTO artists (name)
VALUES ($1)
RETURNING id, name;

-- name: GetArtistById :one
SELECT id, name
FROM artists
WHERE id = $1;

-- name: ListArtists :many
SELECT id, name
FROM artists
WHERE (sqlc.narg('name')::text IS NULL OR name ILIKE '%' || sqlc.narg('name')::text || '%')
ORDER BY id
LIMIT $1 OFFSET $2;


-- name: UpdateArtist :one
UPDATE artists
SET name = $2
WHERE id = $1
RETURNING id, name;

-- name: DeleteArtist :exec
DELETE FROM artists
WHERE id = $1;
