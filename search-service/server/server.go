package server

import (
	"context"
	"log"
	"search-service/config"
	"search-service/document"
	"search-service/proto/search"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	search.UnimplementedSearchServiceServer
	meiliClient *MeiliClient
}

func NewSearchServer(cfg config.MeiliConfig) *SearchServer {
	return &SearchServer{
		meiliClient: NewMeiliClient(cfg),
	}
}

func (s *SearchServer) IndexDocument(ctx context.Context, req *search.IndexRequest) (*search.IndexResponse, error) {
	if req == nil || req.Document == nil {
		return nil, status.Error(codes.InvalidArgument, "Document is required")
	}

	switch doc := req.Document.(type) {
	case *search.IndexRequest_Song:
		log.Println("Indexing Song:", doc.Song.Id)
		err := s.meiliClient.IndexSong(document.SongDocument{
			ID:      doc.Song.Id,
			Title:   doc.Song.Title,
			Artists: doc.Song.Artists,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to index song: %v", err)
		}

	case *search.IndexRequest_Artist:
		log.Println("Indexing Artist:", doc.Artist.Id)
		err := s.meiliClient.IndexArtist(document.ArtistDocument{
			ID:   doc.Artist.Id,
			Name: doc.Artist.Name,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to index artist: %v", err)
		}

	case *search.IndexRequest_Playlist:
		log.Println("Indexing Playlist:", doc.Playlist.Id)
		err := s.meiliClient.IndexPlaylist(document.PlaylistDocument{
			ID:         doc.Playlist.Id,
			SongTitles: doc.Playlist.SongTitles,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to index playlist: %v", err)
		}

	default:
		return nil, status.Error(codes.InvalidArgument, "Unknown document type")
	}

	return &search.IndexResponse{
		Success: true,
		Message: "Document indexed successfully",
	}, nil
}
