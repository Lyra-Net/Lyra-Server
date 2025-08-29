
CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE songs (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    title_token TEXT[],
    categories TEXT[]
);

CREATE TABLE artist_songs (
    artist_id INT NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    song_id TEXT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    PRIMARY KEY (artist_id, song_id)
);

CREATE TABLE playlists (
    playlist_id UUID PRIMARY KEY,
    playlist_name TEXT,
    owner_id UUID NOT NULL REFERENCES users(user_id),
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE playlist_song (
    song_id TEXT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    playlist_id UUID NOT NULL REFERENCES playlists(playlist_id) ON DELETE CASCADE,
    position INT NOT NULL,
    PRIMARY KEY (playlist_id, song_id)
);
