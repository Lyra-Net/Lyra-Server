package dto

import (
	"github.com/google/uuid"
)

type Artist struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CreateSongRequest struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	TitleToken []string `json:"title_token"`
	Categories []string `json:"categories"`
	Duration   int32    `json:"duration"`
	Genre      string   `json:"genre"`
	Mood       string   `json:"mood"`
	Artists    []Artist `json:"artists"`
}
type Playlist struct {
	PlaylistID   uuid.UUID `json:"playlist_id"`
	PlaylistName string    `json:"playlist_name"`
	SongCount    int32     `json:"song_count"`
	SearchKeys   []string  `json:"search_keys"`
}

type Song struct {
	SongID  string   `json:"song_id"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
}
