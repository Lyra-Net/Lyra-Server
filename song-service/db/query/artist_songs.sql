-- name: AddSongArtists :exec
INSERT INTO artist_songs (artist_id, song_id)
VALUES ($1, $2);

-- name: RemoveSongArtists :exec
DELETE FROM artist_songs
WHERE song_id = $1;
