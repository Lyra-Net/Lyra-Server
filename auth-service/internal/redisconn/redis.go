package redisconn

import (
	"auth-service/config"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

type TwoFAEntry struct {
	OTP    string `json:"otp"`
	UserID string `json:"user_id"`
}

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

func Set2faOTP(sessionId, userId, deviceId, otp string) error {
	key := fmt.Sprintf("2fa:%s:%s", sessionId, deviceId)
	entry := TwoFAEntry{
		OTP:    otp,
		UserID: userId,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return Client.Set(Ctx, key, data, 5*time.Minute).Err()
}

func Get2faOTP(sessionId, deviceId string) (otp string, userId string, err error) {
	key := fmt.Sprintf("2fa:%s:%s", sessionId, deviceId)
	data, err := Client.Get(Ctx, key).Bytes()
	if err != nil {
		return "", "", err
	}

	var entry TwoFAEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", "", err
	}
	return entry.OTP, entry.UserID, nil
}

func CheckVerifiedSession(sessionId, deviceId string) bool {
	key := fmt.Sprintf("%s:%s", sessionId, deviceId)
	val, err := Client.Get(Ctx, key).Result()
	if err == redis.Nil || val != "1" {
		return false
	}
	Client.Del(Ctx, key)
	return true
}

func SetVerifiedSession(sessionId, deviceId string) error {
	key := fmt.Sprintf("%s:%s", sessionId, deviceId)
	return Client.Set(Ctx, key, 1, time.Minute*2).Err()
}

func SetVerifing2Fa(userId string) error {
	key := fmt.Sprintf("%s:verifing", userId)
	return Client.Set(Ctx, key, 1, time.Minute*15).Err()
}

func GetVerifing2Fa(userId string) (string, error) {
	key := fmt.Sprintf("%s:verifing", userId)
	return Client.Get(Ctx, key).Result()
}

func RemoveVerifing2Fa(userId string) error {
	key := fmt.Sprintf("%s:verifing", userId)
	return Client.Del(Ctx, key).Err()
}
