package dto

type SongServiceDTO struct {
	CreateArtistRequest CreateArtistRequest `json:"CreateArtistRequest"`
	CreateSongRequest   CreateSongRequest   `json:"CreateSongRequest"`
	CreateSongResponse  CreateSongResponse  `json:"CreateSongResponse"`
	Artist              Artist              `json:"Artist"`
	Song                CreateSongRequest   `json:"Song"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateArtistRequest struct {
	Name string `json:"name"`
}

type CreateSongRequest struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	TitleToken []string `json:"title_token"`
	Categories []string `json:"categories"`
	ArtistIDS  []int64  `json:"artist_ids,omitempty"`
	Artists    []string `json:"artists,omitempty"`
}

type CreateSongResponse struct {
	ID string `json:"id"`
}
