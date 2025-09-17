package server

import (
	"context"
	"log"
	"net"
	"search-service/dto"
	"strings"

	pb "github.com/trandinh0506/BypassBeats/proto/gen/search"
	"google.golang.org/grpc"
)

type searchServer struct {
	pb.UnimplementedSearchServiceServer
	meili MeiliClient
}

func (s *searchServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	q := strings.TrimSpace(req.GetQuery())
	t := strings.ToLower(strings.TrimSpace(req.GetType()))

	const limit int64 = 20
	resp := &pb.SearchResponse{Success: true}

	switch t {
	case "song":
		songsDto, total, err := s.meili.SearchSongs(q, limit)
		if err != nil {
			return &pb.SearchResponse{Success: false, Error: err.Error()}, nil
		}
		resp.Songs = toPBSongs(songsDto)
		resp.TotalHits = int32(total)

	case "artist":
		artistsDto, total, err := s.meili.SearchArtists(q, limit)
		if err != nil {
			return &pb.SearchResponse{Success: false, Error: err.Error()}, nil
		}
		resp.Artists = toPBArtists(artistsDto)
		resp.TotalHits = int32(total)

	case "playlist":
		plsDto, total, err := s.meili.SearchPlaylists(q, limit)
		if err != nil {
			return &pb.SearchResponse{Success: false, Error: err.Error()}, nil
		}
		resp.Playlists = toPBPlaylists(plsDto)
		resp.TotalHits = int32(total)

	default:
		return &pb.SearchResponse{
			Success: false,
			Error:   "unsupported type: " + t,
		}, nil
	}

	return resp, nil
}

func toPBSongs(in []dto.CreateSongRequest) []*pb.Song {
	out := make([]*pb.Song, 0, len(in))
	for _, s := range in {
		artists := make([]*pb.Artist, 0, len(s.Artists))
		for _, a := range s.Artists {
			artists = append(artists, &pb.Artist{
				Id:   a.ID,
				Name: a.Name,
			})
		}
		out = append(out, &pb.Song{
			Id:      s.ID,
			Title:   s.Title,
			Artists: artists,
		})
	}
	return out
}

func toPBArtists(in []dto.Artist) []*pb.Artist {
	out := make([]*pb.Artist, 0, len(in))
	for _, a := range in {
		out = append(out, &pb.Artist{
			Id:   int32(a.ID),
			Name: a.Name,
		})
	}
	return out
}

func toPBPlaylists(in []dto.Playlist) []*pb.Playlist {
	out := make([]*pb.Playlist, 0, len(in))
	for _, p := range in {
		out = append(out, &pb.Playlist{
			PlaylistId:   p.PlaylistID.String(),
			PlaylistName: p.PlaylistName,
			SongCount:    int32(p.SongCount),
			SearchKeys:   p.SearchKeys,
		})
	}
	return out
}

func RunGRPCServer(addr string, meili *MeiliClient) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterSearchServiceServer(s, &searchServer{
		meili: *meili,
	})
	log.Printf("gRPC SearchService listening on %s", addr)
	return s.Serve(lis)
}
