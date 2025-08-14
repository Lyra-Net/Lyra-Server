package server

import (
	"log"
	"search-service/config"
	"search-service/document"

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

func (mc *MeiliClient) IndexSong(song document.SongDocument) error {
	_, err := mc.client.Index("songs").AddDocuments([]document.SongDocument{song}, nil)
	return err
}

func (mc *MeiliClient) IndexArtist(artist document.ArtistDocument) error {
	_, err := mc.client.Index("artists").AddDocuments([]document.ArtistDocument{artist}, nil)
	return err
}

func (mc *MeiliClient) IndexPlaylist(playlist document.PlaylistDocument) error {
	_, err := mc.client.Index("playlists").AddDocuments([]document.PlaylistDocument{playlist}, nil)
	return err
}
