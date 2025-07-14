// internal/services/auth_service.go
package services

import (
	"context"
	"time"

	"identity-service/models"
	"identity-service/proto/auth"
	"identity-service/redisconn"
	"identity-service/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	DB *gorm.DB
}

func NewAuthServer(db *gorm.DB) *AuthServer {
	return &AuthServer{DB: db}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password hashing failed")
	}
	user := models.User{Username: req.Username, Password: string(hashed)}
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already in use")
	}
	return &auth.RegisterResponse{Message: "User created"}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {
	var user models.User
	if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessJti := uuid.New().String()
	refreshJti := uuid.New().String()
	changePassAt := user.UpdatedAt.Unix()
	redisconn.SetChangePassAt(user.ID, changePassAt)

	accessToken, err := utils.GenerateAccessToken(user.ID, accessJti, changePassAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate access token failed")
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, refreshJti)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate refresh token failed")
	}

	tokenRecord := models.RefreshToken{
		ID:        refreshJti,
		UserID:    user.ID,
		Token:     refreshToken,
		DeviceID:  req.DeviceId,
		UserAgent: req.UserAgent,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	s.DB.Create(&tokenRecord)

	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.AuthResponse, error) {
	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	tokenID := claims["jti"].(string)
	userID := uint(claims["user_id"].(float64))
	exp := int64(claims["exp"].(float64))

	if ok, _ := redisconn.IsRefreshTokenBlacklisted(tokenID); ok {
		s.DB.Where("id = ?", userID).Delete(&models.RefreshToken{})
		redisconn.BlacklistAccessToken(tokenID, exp)
		return nil, status.Error(codes.Unauthenticated, "refresh token reuse detected")
	}

	var stored models.RefreshToken
	if err := s.DB.Where("id = ? AND user_id = ?", tokenID, userID).First(&stored).Error; err != nil {
		redisconn.BlacklistRefreshToken(tokenID, exp)
		return nil, status.Error(codes.Unauthenticated, "refresh token invalid or reused")
	}
	s.DB.Delete(&stored)

	newAccessJTI := uuid.New().String()
	newRefreshJTI := uuid.New().String()
	changePassAt, _ := redisconn.GetChangePassAt(userID)

	accessToken, err := utils.GenerateAccessToken(userID, newAccessJTI, changePassAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate access token failed")
	}
	refreshToken, err := utils.GenerateRefreshToken(userID, newRefreshJTI)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate refresh token failed")
	}

	newToken := models.RefreshToken{
		ID:        newRefreshJTI,
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(utils.REFRESH_TOKEN_TIME) * time.Hour),
	}
	s.DB.Create(&newToken)

	return &auth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
