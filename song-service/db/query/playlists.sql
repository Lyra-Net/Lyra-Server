-- name: CreatePlaylist :one
INSERT INTO playlists (playlist_id, playlist_name, owner_id, is_public)
VALUES ($1, $2, $3, $4)
RETURNING playlist_id, playlist_name, owner_id, is_public, created_at, updated_at;

-- name: GetPlaylistById :one
SELECT playlist_id, playlist_name, owner_id, is_public, created_at, updated_at
FROM playlists
WHERE playlist_id = $1;

-- name: ListPlaylists :many
SELECT playlist_id, playlist_name, owner_id, is_public, created_at, updated_at
FROM playlists
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePlaylist :one
UPDATE playlists
SET playlist_name = $2, is_public = $3, updated_at = now()
WHERE playlist_id = $1
RETURNING playlist_id, playlist_name, owner_id, is_public, created_at, updated_at;

-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE playlist_id = $1;
