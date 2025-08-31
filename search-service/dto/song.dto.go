package dto

// Code generated from JSON Schema using quicktype. DO NOT EDIT.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    songServiceDTO, err := UnmarshalSongServiceDTO(bytes)
//    bytes, err = songServiceDTO.Marshal()

import "encoding/json"

func UnmarshalSongServiceDTO(data []byte) (SongServiceDTO, error) {
	var r SongServiceDTO
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SongServiceDTO) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SongServiceDTO struct {
	CreateArtistRequest CreateArtistRequest `json:"CreateArtistRequest"`
	CreateSongRequest   CreateSongRequest   `json:"CreateSongRequest"`
	CreateSongResponse  CreateSongResponse  `json:"CreateSongResponse"`
	Artist              Artist              `json:"Artist"`
	Song                CreateSongRequest   `json:"Song"`
	Playlist            Playlist            `json:"Playlist"`
}

type Artist struct {
	ID   int64  `json:"id"`
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

type Playlist struct {
	PlaylistID   string `json:"playlist_id"`
	PlaylistName string `json:"playlist_name"`
	OwnerID      string `json:"owner_id"`
	IsPublic     bool   `json:"is_public"`
}
