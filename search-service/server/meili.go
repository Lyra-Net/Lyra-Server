package server

import (
	"encoding/json"
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
	}{
		{"songs", "id"},
		{"artists", "id"},
		{"playlists", "id"},
	}

	for _, idx := range indexes {
		_, err := mc.client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        idx.UID,
			PrimaryKey: idx.PrimaryKey,
		})
		if err != nil {
			log.Printf("Index %s maybe existed: %v", idx.UID, err)
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

func (mc *MeiliClient) SearchSongs(query string, limit int64) ([]dto.CreateSongRequest, error) {
	res, err := mc.client.Index("songs").Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res.Hits)
	if err != nil {
		return nil, err
	}
	var songs []dto.CreateSongRequest
	if err := json.Unmarshal(b, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func (mc *MeiliClient) SearchArtists(query string, limit int64) ([]dto.Artist, error) {
	res, err := mc.client.Index("artists").Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res.Hits)
	if err != nil {
		return nil, err
	}
	var artists []dto.Artist
	if err := json.Unmarshal(b, &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func (mc *MeiliClient) SearchPlaylists(query string, limit int64) ([]dto.PlaylistDTO, error) {
	res, err := mc.client.Index("playlists").Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res.Hits)
	if err != nil {
		return nil, err
	}
	var pls []dto.PlaylistDTO
	if err := json.Unmarshal(b, &pls); err != nil {
		return nil, err
	}
	return pls, nil
}
