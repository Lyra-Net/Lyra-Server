package playlist

import (
	"context"
	"song-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/trandinh0506/BypassBeats/proto/gen/playlist"
)

type PlaylistService struct {
	pb.UnimplementedPlaylistServiceServer
	q *repository.Queries
}

func NewPlaylistService(q *repository.Queries) *PlaylistService {
	return &PlaylistService{q: q}
}

// ================== CreatePlaylist ==================

func (s *PlaylistService) CreatePlaylist(ctx context.Context, req *pb.CreatePlaylistRequest) (*pb.CreatePlaylistResponse, error) {
	playlist, err := s.q.CreatePlaylist(ctx, repository.CreatePlaylistParams{
		PlaylistID:   uuid.New(),
		PlaylistName: pgtype.Text{String: req.PlaylistName, Valid: req.PlaylistName != ""},
		OwnerID:      uuid.MustParse(req.OwnerId),
		IsPublic:     pgtype.Bool{Bool: req.IsPublic, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreatePlaylistResponse{
		PlaylistId: playlist.PlaylistID.String(),
	}, nil
}

// ================== GetPlaylistByID ==================

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, req *pb.GetPlaylistByIDRequest) (*pb.Playlist, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, err
	}

	playlist, err := s.q.GetPlaylistById(ctx, playlistUUID)
	if err != nil {
		return nil, err
	}

	return &pb.Playlist{
		PlaylistId:   playlist.PlaylistID.String(),
		PlaylistName: playlist.PlaylistName.String,
		OwnerId:      playlist.OwnerID.String(),
		IsPublic:     playlist.IsPublic.Bool,
	}, nil
}

// ================== ListMyPlaylists ==================

func (s *PlaylistService) ListMyPlaylists(ctx context.Context, req *pb.ListMyPlaylistsRequest) (*pb.ListMyPlaylistsResponse, error) {
	ownerUUID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, err
	}

	playlists, err := s.q.ListMyPlaylists(ctx, ownerUUID) // cần viết thêm SQL query
	if err != nil {
		return nil, err
	}

	resp := &pb.ListMyPlaylistsResponse{}
	for _, p := range playlists {
		resp.Playlists = append(resp.Playlists, &pb.Playlist{
			PlaylistId:   p.PlaylistID.String(),
			PlaylistName: p.PlaylistName.String,
			OwnerId:      p.OwnerID.String(),
			IsPublic:     p.IsPublic.Bool,
		})
	}
	return resp, nil
}

// ================== RemoveSongFromPlaylist ==================

func (s *PlaylistService) RemoveSongFromPlaylist(ctx context.Context, req *pb.RemoveSongFromPlaylistRequest) (*emptypb.Empty, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, err
	}

	err = s.q.RemoveSongFromPlaylist(ctx, repository.RemoveSongFromPlaylistParams{
		PlaylistID: playlistUUID,
		SongID:     req.SongId,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ================== MoveSongInPlaylist ==================

func (s *PlaylistService) MoveSongInPlaylist(ctx context.Context, req *pb.MoveSongInPlaylistRequest) (*emptypb.Empty, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, err
	}

	oldPos, err := s.q.GetSongPosition(ctx, repository.GetSongPositionParams{
		PlaylistID: playlistUUID,
		SongID:     req.SongId,
	})
	if err != nil {
		return nil, err
	}

	newPos := req.NewPosition

	// Transaction
	_, pgxTx, err := s.q.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	tx := s.q.WithTx(pgxTx)
	defer pgxTx.Rollback(ctx)

	if newPos > oldPos {
		if err := tx.ShiftPositionsDown(ctx, repository.ShiftPositionsDownParams{
			PlaylistID: playlistUUID,
			Position:   oldPos,
			Position_2: newPos,
		}); err != nil {
			return nil, err
		}
	} else if newPos < oldPos {
		if err := tx.ShiftPositionsUp(ctx, repository.ShiftPositionsUpParams{
			PlaylistID: playlistUUID,
			Position:   newPos,
			Position_2: oldPos,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.UpdateSongPosition(ctx, repository.UpdateSongPositionParams{
		PlaylistID: playlistUUID,
		SongID:     req.SongId,
		Position:   newPos,
	}); err != nil {
		return nil, err
	}

	if err := pgxTx.Commit(ctx); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
