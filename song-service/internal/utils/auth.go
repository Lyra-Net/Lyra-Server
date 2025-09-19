package utils

import (
	"context"
	"song-service/internal/interceptor"
	"song-service/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(interceptor.UserIDKey)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

func CheckPlaylistOwner(ctx context.Context, q *repository.Queries, playlistID uuid.UUID) (uuid.UUID, error) {
	userId, ok := GetUserID(ctx)
	if !ok {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	ownerUUID, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	playlist, err := q.GetPlaylistById(ctx, playlistID)
	if err != nil {
		return uuid.Nil, status.Errorf(codes.NotFound, "playlist not found")
	}

	if playlist.OwnerID != ownerUUID {
		return uuid.Nil, status.Errorf(codes.PermissionDenied, "you do not own this playlist")
	}

	return ownerUUID, nil
}
