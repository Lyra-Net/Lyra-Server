package services

import (
	"context"
	db "identity-service/internal/repository"
	"identity-service/proto/auth"
	"identity-service/redisconn"
	"identity-service/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Queries *db.Queries
}

func NewAuthServer(q *db.Queries) *AuthServer {
	return &AuthServer{Queries: q}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password hashing failed")
	}

	userID := uuid.New()
	err = s.Queries.CreateUser(ctx, db.CreateUserParams{
		UserID:       userID,
		Username:     req.Username,
		PasswordHash: string(hashed),
	})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already in use")
	}

	return &auth.RegisterResponse{Message: "User created"}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {
	user, err := s.Queries.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessJti := uuid.New()
	refreshJti := uuid.New()
	changePassAt := user.UpdatedAt.Time.Unix()
	redisconn.SetChangePassAt(user.UserID, changePassAt)

	accessToken, err := utils.GenerateAccessToken(user.UserID, accessJti, changePassAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate access token failed")
	}
	refreshToken, err := utils.GenerateRefreshToken(user.UserID, refreshJti)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate refresh token failed")
	}

	err = s.Queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        refreshJti,
		UserID:    user.UserID,
		Token:     refreshToken,
		DeviceID:  pgtype.Text{String: req.DeviceId, Valid: req.DeviceId != ""},
		UserAgent: pgtype.Text{String: req.UserAgent, Valid: req.UserAgent != ""},
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "store refresh token failed")
	}

	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
