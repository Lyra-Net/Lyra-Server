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

-- name: GetPlaylistWithSongs :many
SELECT 
    p.playlist_id,
    p.playlist_name,
    p.owner_id,
    p.is_public,
    p.created_at,
    p.updated_at,
    ps.song_id,
    ps.position,
    s.title,
    s.title_token,
    s.categories
FROM playlists p
LEFT JOIN playlist_song ps ON p.playlist_id = ps.playlist_id
LEFT JOIN songs s ON ps.song_id = s.id
WHERE p.playlist_id = $1
ORDER BY ps.position ASC;

-- name: GetSongsInPlaylist :many
SELECT 
    ps.song_id,
    ps.position,
    s.title,
    s.title_token,
    s.categories
FROM playlist_song ps
JOIN songs s ON ps.song_id = s.id
WHERE ps.playlist_id = $1
ORDER BY ps.position ASC;

-- name: ListMyPlaylists :many
SELECT playlist_id, playlist_name, owner_id, is_public, created_at, updated_at
FROM playlists
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: GetPlaylistWithSongsAndArtists :many
SELECT 
    p.playlist_id,
    p.playlist_name,
    p.owner_id,
    p.is_public,
    p.created_at,
    p.updated_at,
    ps.song_id,
    ps.position,
    s.title,
    a.id AS artist_id,
    a.name AS artist_name
FROM playlists p
LEFT JOIN playlist_song ps ON p.playlist_id = ps.playlist_id
LEFT JOIN songs s ON ps.song_id = s.id
LEFT JOIN artist_songs sa ON s.id = sa.song_id
LEFT JOIN artists a ON sa.artist_id = a.id
WHERE p.playlist_id = $1
ORDER BY ps.position ASC, a.name ASC;


-- name: ListMyPlaylistsWithSongsAndArtists :many
SELECT 
    p.playlist_id,
    p.playlist_name,
    p.owner_id,
    p.is_public,
    ps.song_id,
    ps.position,
    s.title,
    a.id AS artist_id,
    a.name AS artist_name
FROM playlists p
LEFT JOIN playlist_song ps ON p.playlist_id = ps.playlist_id
LEFT JOIN songs s ON ps.song_id = s.id
LEFT JOIN artist_songs sa ON s.id = sa.song_id
LEFT JOIN artists a ON sa.artist_id = a.id
WHERE p.owner_id = $1
ORDER BY p.created_at DESC, ps.position ASC, a.name ASC;
