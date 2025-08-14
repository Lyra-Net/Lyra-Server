package document

type PlaylistDocument struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	OwnerID    *string  `json:"owner_id,omitempty"`
	OwnerName  *string  `json:"owner_name,omitempty"`
	SongTitles []string `json:"song_titles"`
	Tags       []string `json:"tags"`
	IsPublic   bool     `json:"is_public"`
}
