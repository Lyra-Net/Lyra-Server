package server

import (
	"encoding/json"
	"fmt"
	"log"
	"search-service/config"
	"search-service/dto"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliClient struct {
	client meilisearch.ServiceManager
}

func NewMeiliClient(cfg config.MeiliConfig) *MeiliClient {
	client := meilisearch.New(cfg.Host, meilisearch.WithAPIKey(cfg.APIKey))
	mc := &MeiliClient{client: client}
	mc.initIndexes()
	return mc
}

func (mc *MeiliClient) initIndexes() {
	indexes := []struct {
		UID        string
		PrimaryKey string
		Searchable []string
		Filterable []interface{}
	}{
		{
			"songs", "id",
			[]string{"title", "artists.name", "categories"},
			[]interface{}{"id"},
		},
		{
			"artists", "id",
			[]string{"name"},
			[]interface{}{"id"},
		},
		{
			"playlists", "playlist_id",
			[]string{"playlist_name", "search_keys"},
			[]interface{}{"playlist_id"},
		},
	}

	for _, idx := range indexes {
		_, err := mc.client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        idx.UID,
			PrimaryKey: idx.PrimaryKey,
		})
		if err != nil {
			log.Printf("Index %s maybe existed: %v", idx.UID, err)
		}

		if _, err := mc.client.Index(idx.UID).UpdateSearchableAttributes(&idx.Searchable); err != nil {
			log.Printf("Set searchable attrs for %s failed: %v", idx.UID, err)
		}
		if _, err := mc.client.Index(idx.UID).UpdateFilterableAttributes(&(idx.Filterable)); err != nil {
			log.Printf("Set filterable attrs for %s failed: %v", idx.UID, err)
		}
	}
}

func (mc *MeiliClient) IndexSong(req dto.CreateSongRequest) error {
	_, err := mc.client.Index("songs").AddDocuments([]dto.CreateSongRequest{req}, nil)
	return err
}

func (mc *MeiliClient) IndexArtist(req dto.Artist) error {
	_, err := mc.client.Index("artists").AddDocuments([]dto.Artist{req}, nil)
	return err
}

func (mc *MeiliClient) IndexPlaylist(req dto.Playlist) error {
	_, err := mc.client.Index("playlists").AddDocuments([]dto.Playlist{req}, nil)
	return err
}

func (mc *MeiliClient) SearchSongs(query string, limit int64) ([]dto.CreateSongRequest, int64, error) {
	res, err := mc.client.Index("songs").Search(query, &meilisearch.SearchRequest{Limit: limit})
	if err != nil {
		return nil, 0, err
	}
	b, _ := json.Marshal(res.Hits)
	var songs []dto.CreateSongRequest
	if err := json.Unmarshal(b, &songs); err != nil {
		return nil, 0, err
	}
	return songs, res.EstimatedTotalHits, nil
}

func (mc *MeiliClient) SearchArtists(query string, limit int64) ([]dto.Artist, int64, error) {
	res, err := mc.client.Index("artists").Search(query, &meilisearch.SearchRequest{Limit: limit})
	if err != nil {
		return nil, 0, err
	}
	b, _ := json.Marshal(res.Hits)
	var artists []dto.Artist
	if err := json.Unmarshal(b, &artists); err != nil {
		return nil, 0, err
	}
	return artists, res.EstimatedTotalHits, nil
}

func (mc *MeiliClient) SearchPlaylists(query string, limit int64) ([]dto.Playlist, int64, error) {
	res, err := mc.client.Index("playlists").Search(query, &meilisearch.SearchRequest{Limit: limit})
	if err != nil {
		return nil, 0, err
	}
	b, _ := json.Marshal(res.Hits)
	var pls []dto.Playlist
	if err := json.Unmarshal(b, &pls); err != nil {
		return nil, 0, err
	}
	return pls, res.EstimatedTotalHits, nil
}

// Add document with same ID will replace the old one <=> update

func (mc *MeiliClient) UpdateSong(req dto.CreateSongRequest) error {
	_, err := mc.client.Index("songs").AddDocuments([]dto.CreateSongRequest{req}, nil)
	return err
}

func (mc *MeiliClient) DeleteSong(id string) error {
	_, err := mc.client.Index("songs").DeleteDocument(id)
	return err
}

func (mc *MeiliClient) UpdateArtist(req dto.Artist) error {
	_, err := mc.client.Index("artists").AddDocuments([]dto.Artist{req}, nil)
	return err
}

func (mc *MeiliClient) DeleteArtist(id int32) error {
	_, err := mc.client.Index("artists").DeleteDocument(fmt.Sprint(id))
	return err
}

func (mc *MeiliClient) UpdatePlaylist(req dto.Playlist) error {
	_, err := mc.client.Index("playlists").AddDocuments([]dto.Playlist{req}, nil)
	return err
}

func (mc *MeiliClient) DeletePlaylist(id string) error {
	_, err := mc.client.Index("playlists").DeleteDocument(id)
	return err
}
