-- name: AddSongToPlaylist :exec
INSERT INTO playlist_song (playlist_id, song_id, position)
VALUES ($1, $2, $3);

-- name: RemoveSongFromPlaylist :exec
DELETE FROM playlist_song
WHERE playlist_id = $1 AND song_id = $2;
