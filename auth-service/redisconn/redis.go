package redisconn

import (
	"context"
	"fmt"
	"identity-service/config"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedis() {
	cfg := config.GetConfig()
	redisURL := cfg.REDIS_URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_URL: %v", err))
	}
	Client = redis.NewClient(opt)
}

// SetChangePassAt caches the password change timestamp (unix) for a given user
func SetChangePassAt(userID uuid.UUID, ts int64) error {
	key := fmt.Sprintf("changepass:user:%d", userID)
	return Client.Set(Ctx, key, ts, time.Hour*1).Err() // optional TTL
}

// GetChangePassAt retrieves the cached password change timestamp
func GetChangePassAt(userID uuid.UUID) (int64, error) {
	key := fmt.Sprintf("changepass:user:%d", userID)
	val, err := Client.Get(Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// DeleteChangePassAt removes the cached changePassAt for a user (optional)
func DeleteChangePassAt(userID uuid.UUID) error {
	key := fmt.Sprintf("changepass:user:%d", userID)
	return Client.Del(Ctx, key).Err()
}

// BlacklistAccessToken stores a revoked access token's jti with expiration
func BlacklistAccessToken(jti uuid.UUID, exp int64) error {
	key := fmt.Sprintf("blacklist:access:%s", jti.String())
	expiry := time.Until(time.Unix(exp, 0))
	return Client.Set(Ctx, key, "revoked", expiry).Err()
}

func IsAccessTokenBlacklisted(jti uuid.UUID) (bool, error) {
	key := fmt.Sprintf("blacklist:access:%s", jti.String())
	_, err := Client.Get(Ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}

func BlacklistRefreshToken(jti uuid.UUID, exp int64) error {
	key := fmt.Sprintf("blacklist:refresh:%s", jti.String())
	expiry := time.Until(time.Unix(exp, 0))
	return Client.Set(Ctx, key, "revoked", expiry).Err()
}

func IsRefreshTokenBlacklisted(jti uuid.UUID) (bool, error) {
	key := fmt.Sprintf("blacklist:refresh:%s", jti.String())
	_, err := Client.Get(Ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}
