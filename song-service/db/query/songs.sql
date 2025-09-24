-- name: CreateSong :one
INSERT INTO songs (id, title, title_token, categories, duration, genre, mood)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, title_token, categories, duration, genre, mood;

-- name: GetSongById :one
SELECT id, title, title_token, categories, duration, genre, mood
FROM songs
WHERE id = $1;

-- name: ListSongs :many
SELECT id, title, title_token, categories, duration, genre, mood
FROM songs
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateSong :one
UPDATE songs
SET title       = $2,
    title_token = $3,
    categories  = $4,
    duration    = $5,
    genre       = $6,
    mood        = $7
WHERE id = $1
RETURNING id, title, title_token, categories, duration, genre, mood;

-- name: DeleteSong :exec
DELETE FROM songs
WHERE id = $1;

-- name: ListSongsWithArtists :many
SELECT 
    s.id, 
    s.title, 
    s.title_token, 
    s.categories,
    s.duration,
    s.genre,
    s.mood,
    COALESCE(
        json_agg(json_build_object('id', a.id, 'name', a.name)) 
        FILTER (WHERE a.id IS NOT NULL),
        '[]'
    ) AS artists
FROM 
    songs s
LEFT JOIN 
    artist_songs sa ON sa.song_id = s.id
LEFT JOIN 
    artists a ON a.id = sa.artist_id
GROUP BY 
    s.id, s.title, s.title_token, s.categories, s.duration, s.genre, s.mood
ORDER BY 
    s.id
LIMIT $1 OFFSET $2;
