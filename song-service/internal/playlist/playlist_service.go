package playlist

import (
	"context"
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
		return nil, status.Errorf(codes.InvalidArgument, "invalid playlist id")
	}

	rows, err := s.q.GetPlaylistWithSongsAndArtists(ctx, playlistUUID)
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
		ownerUUID, _ := uuid.Parse(userId)
		if rows[0].OwnerID != ownerUUID {
			return nil, status.Errorf(codes.PermissionDenied, "you do not have access to this playlist")
		}
	}

	res := &pb.Playlist{
		PlaylistId:   rows[0].PlaylistID.String(),
		PlaylistName: rows[0].PlaylistName.String,
		OwnerId:      rows[0].OwnerID.String(),
		IsPublic:     rows[0].IsPublic.Bool,
		Songs:        []*pb.PlaylistSong{},
	}

	songMap := make(map[string]*pb.PlaylistSong)

	for _, row := range rows {
		if !row.SongID.Valid {
			continue
		}

		songID := row.SongID.String
		song, exists := songMap[songID]
		if !exists {
			song = &pb.PlaylistSong{
				SongId:   songID,
				Title:    row.Title.String,
				Position: int32(row.Position.Int32),
				Artists:  []*pb.Artist{},
			}
			songMap[songID] = song
			res.Songs = append(res.Songs, song)
		}

		if row.ArtistID.Valid {
			song.Artists = append(song.Artists, &pb.Artist{
				Id:   int32(row.ArtistID.Int32),
				Name: row.ArtistName.String,
			})
		}
	}

	return res, nil
}

// ================== ListMyPlaylists ==================

func (s *PlaylistService) ListMyPlaylists(ctx context.Context, req *pb.ListMyPlaylistsRequest) (*pb.ListMyPlaylistsResponse, error) {
	userId, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	ownerUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	rows, err := s.q.ListMyPlaylistsWithSongsAndArtists(ctx, ownerUUID)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListMyPlaylistsResponse{Playlists: []*pb.Playlist{}}
	playlistMap := make(map[string]*pb.Playlist)

	for _, row := range rows {
		pl, exists := playlistMap[row.PlaylistID.String()]
		if !exists {
			pl = &pb.Playlist{
				PlaylistId:   row.PlaylistID.String(),
				PlaylistName: row.PlaylistName.String,
				OwnerId:      row.OwnerID.String(),
				IsPublic:     row.IsPublic.Bool,
				Songs:        []*pb.PlaylistSong{},
			}
			playlistMap[row.PlaylistID.String()] = pl
			resp.Playlists = append(resp.Playlists, pl)
		}

		if row.SongID.Valid {
			songID := row.SongID.String
			var song *pb.PlaylistSong
			for _, sng := range pl.Songs {
				if sng.SongId == songID {
					song = sng
					break
				}
			}
			if song == nil {
				song = &pb.PlaylistSong{
					SongId:   songID,
					Title:    row.Title.String,
					Position: int32(row.Position.Int32),
					Artists:  []*pb.Artist{},
				}
				pl.Songs = append(pl.Songs, song)
			}

			if row.ArtistID.Valid {
				song.Artists = append(song.Artists, &pb.Artist{
					Id:   int32(row.ArtistID.Int32),
					Name: row.ArtistName.String,
				})
			}
		}
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
