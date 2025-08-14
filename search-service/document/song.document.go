package document

type SongDocument struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Artists   []string `json:"artists"`
	CreatedAt int64    `json:"created_at"`
}
