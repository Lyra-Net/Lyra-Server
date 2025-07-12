package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret       = []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET"))
	refreshSecret      = []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET"))
	ACCESS_TOKEN_TIME  = 1
	REFRESH_TOKEN_TIME = 7 * 24
)

// ================= ACCESS TOKEN =====================

func GenerateAccessToken(userID uint, jti string, changePassAt int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":        userID,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(time.Duration(ACCESS_TOKEN_TIME) * time.Hour).Unix(),
		"jti":            jti,
		"change_pass_at": changePassAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid access token claims")
	}
	return claims, nil
}

// ================= REFRESH TOKEN =====================

func GenerateRefreshToken(userID uint, jti string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Duration(REFRESH_TOKEN_TIME) * time.Hour).Unix(),
		"jti":     jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ParseRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}
	return claims, nil
}
