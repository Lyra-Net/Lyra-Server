package playlist

import (
	"context"
	"log"
	"song-service/internal/repository"
	"song-service/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	userId, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	playlist, err := s.q.CreatePlaylist(ctx, repository.CreatePlaylistParams{
		PlaylistID:   uuid.New(),
		PlaylistName: pgtype.Text{String: req.PlaylistName, Valid: req.PlaylistName != ""},
		OwnerID:      uuid.MustParse(userId),
		IsPublic:     pgtype.Bool{Bool: req.IsPublic, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreatePlaylistResponse{
		PlaylistId: playlist.PlaylistID.String(),
		IsSuccess:  true,
	}, nil
}

// ================== GetPlaylistByID ==================

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, req *pb.GetPlaylistByIDRequest) (*pb.Playlist, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, err
	}

	rows, err := s.q.GetPlaylistWithSongs(ctx, playlistUUID)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, status.Errorf(codes.NotFound, "playlist not found")
	}
	if !rows[0].IsPublic.Bool {
		userId, ok := utils.GetUserID(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
		}
		ownerUUID, err := uuid.Parse(userId)
		if err != nil {
			return nil, err
		}
		if rows[0].OwnerID != ownerUUID {
			return nil, status.Errorf(codes.PermissionDenied, "you do not have access to this playlist")
		}
	}

	res := &pb.Playlist{
		PlaylistId:   rows[0].PlaylistID.String(),
		PlaylistName: rows[0].PlaylistName.String,
		OwnerId:      rows[0].OwnerID.String(),
		IsPublic:     rows[0].IsPublic.Bool,
		Songs:        make([]*pb.PlaylistSong, 0),
	}

	for _, row := range rows {
		if row.SongID.Valid {
			res.Songs = append(res.Songs, &pb.PlaylistSong{
				SongId:     row.SongID.String,
				Title:      row.Title.String,
				TitleToken: row.TitleToken,
				Categories: row.Categories,
				Position:   int32(row.Position.Int32),
			})
		}
	}

	return res, nil
}

// ================== ListMyPlaylists ==================

func (s *PlaylistService) ListMyPlaylists(ctx context.Context, req *pb.ListMyPlaylistsRequest) (*pb.ListMyPlaylistsResponse, error) {
	userId, ok := utils.GetUserID(ctx)
	log.Println("UserID:", userId, "OK:", ok)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	ownerUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	playlists, err := s.q.ListMyPlaylists(ctx, ownerUUID)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListMyPlaylistsResponse{}
	for _, p := range playlists {
		pl := &pb.Playlist{
			PlaylistId:   p.PlaylistID.String(),
			PlaylistName: p.PlaylistName.String,
			OwnerId:      p.OwnerID.String(),
			IsPublic:     p.IsPublic.Bool,
			Songs:        make([]*pb.PlaylistSong, 0),
		}

		songs, err := s.q.GetSongsInPlaylist(ctx, p.PlaylistID)
		if err != nil {
			return nil, err
		}

		for _, sng := range songs {
			pl.Songs = append(pl.Songs, &pb.PlaylistSong{
				SongId:     sng.SongID,
				Title:      sng.Title,
				TitleToken: sng.TitleToken,
				Categories: sng.Categories,
				Position:   sng.Position,
			})
		}

		resp.Playlists = append(resp.Playlists, pl)
	}

	return resp, nil
}

// ================== AddSongToPlaylist ==================

func (s *PlaylistService) AddSongToPlaylist(ctx context.Context, req *pb.AddSongToPlaylistRequest) (*pb.AddSongToPlaylistResponse, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid playlist id")
	}

	_, err = utils.CheckPlaylistOwner(ctx, s.q, playlistUUID)
	if err != nil {
		return nil, err
	}

	err = s.q.AddSongToPlaylist(ctx, repository.AddSongToPlaylistParams{
		PlaylistID: playlistUUID,
		SongID:     req.SongId,
	})
	if err != nil {
		return &pb.AddSongToPlaylistResponse{IsSuccess: false}, err
	}

	return &pb.AddSongToPlaylistResponse{IsSuccess: true}, nil
}

// ================== RemoveSongFromPlaylist ==================

func (s *PlaylistService) RemoveSongFromPlaylist(ctx context.Context, req *pb.RemoveSongFromPlaylistRequest) (*pb.RemoveSongFromPlaylistResponse, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid playlist id")
	}

	_, err = utils.CheckPlaylistOwner(ctx, s.q, playlistUUID)
	if err != nil {
		return nil, err
	}

	err = s.q.RemoveSongFromPlaylist(ctx, repository.RemoveSongFromPlaylistParams{
		PlaylistID: playlistUUID,
		SongID:     req.SongId,
	})
	if err != nil {
		return &pb.RemoveSongFromPlaylistResponse{IsSuccess: false}, err
	}

	return &pb.RemoveSongFromPlaylistResponse{IsSuccess: true}, nil
}

// ================== MoveSongInPlaylist ==================

func (s *PlaylistService) MoveSongInPlaylist(ctx context.Context, req *pb.MoveSongInPlaylistRequest) (*pb.MoveSongInPlaylistResponse, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid playlist id")
	}

	_, err = utils.CheckPlaylistOwner(ctx, s.q, playlistUUID)
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

	return &pb.MoveSongInPlaylistResponse{IsSuccess: true}, nil
}

// ================== UpdatePlaylist ==================

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, req *pb.UpdatePlaylistRequest) (*pb.UpdatePlaylistResponse, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid playlist id")
	}

	_, err = utils.CheckPlaylistOwner(ctx, s.q, playlistUUID)
	if err != nil {
		return nil, err
	}

	pl, err := s.q.UpdatePlaylist(ctx, repository.UpdatePlaylistParams{
		PlaylistID:   playlistUUID,
		PlaylistName: pgtype.Text{String: req.PlaylistName, Valid: req.PlaylistName != ""},
		IsPublic:     pgtype.Bool{Bool: req.IsPublic, Valid: true},
	})
	if err != nil {
		return &pb.UpdatePlaylistResponse{IsSuccess: false}, err
	}

	return &pb.UpdatePlaylistResponse{
		Playlist: &pb.Playlist{
			PlaylistId:   pl.PlaylistID.String(),
			PlaylistName: pl.PlaylistName.String,
			OwnerId:      pl.OwnerID.String(),
			IsPublic:     pl.IsPublic.Bool,
		},
		IsSuccess: true,
	}, nil
}

// ================== DeletePlaylist ==================

func (s *PlaylistService) DeletePlaylist(ctx context.Context, req *pb.DeletePlaylistRequest) (*pb.DeletePlaylistResponse, error) {
	playlistUUID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid playlist id")
	}

	_, err = utils.CheckPlaylistOwner(ctx, s.q, playlistUUID)
	if err != nil {
		return nil, err
	}

	err = s.q.DeletePlaylist(ctx, playlistUUID)
	if err != nil {
		return &pb.DeletePlaylistResponse{IsSuccess: false}, err
	}

	return &pb.DeletePlaylistResponse{IsSuccess: true}, nil
}
