package server

import (
	"context"
	"log"
	"net"
	"search-service/dto"
	"strings"

	pb "github.com/trandinh0506/BypassBeats/proto/gen/search/proto"
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
		songsDto, err := s.meili.SearchSongs(q, limit)
		if err != nil {
			return &pb.SearchResponse{Success: false, Error: err.Error()}, nil
		}
		resp.Songs = toPBSongs(songsDto)

	case "artist":
		artistsDto, err := s.meili.SearchArtists(q, limit)
		if err != nil {
			return &pb.SearchResponse{Success: false, Error: err.Error()}, nil
		}
		resp.Artists = toPBArtists(artistsDto)

	case "playlist":
		plsDto, err := s.meili.SearchPlaylists(q, limit)
		if err != nil {
			return &pb.SearchResponse{Success: false, Error: err.Error()}, nil
		}
		resp.Playlists = toPBPlaylists(plsDto)

	// nếu client không truyền type → trả về cả 3
	case "", "all":
		if songsDto, err := s.meili.SearchSongs(q, limit); err == nil {
			resp.Songs = toPBSongs(songsDto)
		} else {
			resp.Success = false
			resp.Error = err.Error()
			return resp, nil
		}
		if artistsDto, err := s.meili.SearchArtists(q, limit); err == nil {
			resp.Artists = toPBArtists(artistsDto)
		} else {
			resp.Success = false
			resp.Error = err.Error()
			return resp, nil
		}
		if plsDto, err := s.meili.SearchPlaylists(q, limit); err == nil {
			resp.Playlists = toPBPlaylists(plsDto)
		} else {
			resp.Success = false
			resp.Error = err.Error()
			return resp, nil
		}

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
		out = append(out, &pb.Song{
			Id:    s.ID,
			Title: s.Title,
			// Chú ý: proto có field Artist string,
			// nhưng DTO của bạn đang là ArtistIDS []int.
			// Hoặc đổi proto, hoặc index thêm field artist name.
			Artist: "", // TODO: điền nếu bạn index sẵn tên artist
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
			PlaylistId:   p.PlaylistID,
			PlaylistName: p.PlaylistName,
			OwnerId:      p.OwnerID,
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
