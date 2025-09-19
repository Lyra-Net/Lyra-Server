-- name: AddSongToPlaylist :exec
INSERT INTO playlist_song (playlist_id, song_id, position)
VALUES (
    $1,
    $2,
    COALESCE((
        SELECT MAX(position) + 1
        FROM playlist_song
        WHERE playlist_id = $1
    ), 1)
);

-- name: RemoveSongFromPlaylist :exec
DELETE FROM playlist_song
WHERE playlist_id = $1 AND song_id = $2;

-- name: ClearPlaylist :exec
DELETE FROM playlist_song
WHERE playlist_id = $1;

-- name: GetSongPosition :one
SELECT position
FROM playlist_song
WHERE playlist_id = $1 AND song_id = $2;

-- name: ShiftPositionsDown :exec
UPDATE playlist_song
SET position = position - 1
WHERE playlist_id = $1
  AND position > $2 AND position <= $3;

-- name: ShiftPositionsUp :exec
UPDATE playlist_song
SET position = position + 1
WHERE playlist_id = $1
  AND position >= $2 AND position < $3;

-- name: UpdateSongPosition :exec
UPDATE playlist_song
SET position = $3
WHERE playlist_id = $1 AND song_id = $2;
