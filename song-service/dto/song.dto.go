package dto

type Artist struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CreateArtistRequest struct {
	Name string `json:"name"`
}

type CreatePlaylistRequest struct {
	PlaylistName string  `json:"playlist_name"`
	OwnerID      string  `json:"owner_id"`
	IsPublic     bool    `json:"is_public"`
	PlaylistID   *string `json:"playlist_id,omitempty"`
}

type CreatePlaylistResponse struct {
	PlaylistID string `json:"playlist_id"`
}

type CreateSongRequest struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	TitleToken []string `json:"title_token"`
	Categories []string `json:"categories"`
	Artists    []Artist `json:"artists"`
}

type CreateSongResponse struct {
	ID string `json:"id"`
}

type GetPlaylistByIDResponse struct {
	PlaylistID   string         `json:"playlist_id"`
	PlaylistName string         `json:"playlist_name"`
	Songs        []SongResponse `json:"songs"`
}

type SongResponse struct {
	SongID     string   `json:"song_id"`
	Title      string   `json:"title"`
	Categories []string `json:"categories"`
	Position   int32    `json:"position"`
}

type MoveSongRequest struct {
	SongID      string `json:"song_id" binding:"required"`
	NewPosition int32  `json:"new_position" binding:"required"`
}

type MoveSongResponse struct {
	PlaylistID  string `json:"playlist_id"`
	SongID      string `json:"song_id"`
	NewPosition int32  `json:"new_position"`
}

type PlaylistDTO struct {
	PlaylistID   string `json:"playlist_id"`
	PlaylistName string `json:"playlist_name"`
	IsPublic     bool   `json:"is_public"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
