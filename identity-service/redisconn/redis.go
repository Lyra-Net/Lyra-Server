package redisconn

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_URL: %v", err))
	}
	Client = redis.NewClient(opt)
}

// SetChangePassAt caches the password change timestamp (unix) for a given user
func SetChangePassAt(userID uint, ts int64) error {
	key := fmt.Sprintf("changepass:user:%d", userID)
	return Client.Set(Ctx, key, ts, time.Hour*1).Err() // optional TTL
}

// GetChangePassAt retrieves the cached password change timestamp
func GetChangePassAt(userID uint) (int64, error) {
	key := fmt.Sprintf("changepass:user:%d", userID)
	val, err := Client.Get(Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// DeleteChangePassAt removes the cached changePassAt for a user (optional)
func DeleteChangePassAt(userID uint) error {
	key := fmt.Sprintf("changepass:user:%d", userID)
	return Client.Del(Ctx, key).Err()
}

// BlacklistAccessToken stores a revoked access token's jti with expiration
func BlacklistAccessToken(jti string, exp int64) error {
	key := fmt.Sprintf("blacklist:access:%s", jti)
	expiry := time.Until(time.Unix(exp, 0))
	return Client.Set(Ctx, key, "revoked", expiry).Err()
}

// IsAccessTokenBlacklisted checks if access token jti is blacklisted
func IsAccessTokenBlacklisted(jti string) (bool, error) {
	key := fmt.Sprintf("blacklist:access:%s", jti)
	_, err := Client.Get(Ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}

// BlacklistRefreshToken stores a revoked refresh token id with expiration
func BlacklistRefreshToken(tokenID string, exp int64) error {
	key := fmt.Sprintf("blacklist:refresh:%s", tokenID)
	expiry := time.Until(time.Unix(exp, 0))
	return Client.Set(Ctx, key, "revoked", expiry).Err()
}

// IsRefreshTokenBlacklisted checks if refresh token id is blacklisted
func IsRefreshTokenBlacklisted(tokenID string) (bool, error) {
	key := fmt.Sprintf("blacklist:refresh:%s", tokenID)
	_, err := Client.Get(Ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}
