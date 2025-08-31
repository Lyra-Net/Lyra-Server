package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type MeiliConfig struct {
	Host   string
	APIKey string
}

type AppConfig struct {
	Meili    MeiliConfig
	GRPCPort string
	Brokers  []string
	Topic    string
	GroupID  string
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env vars")
	}
	brokers := strings.Split(getEnv("KAFKA_BROKER", ""), ",")
	return &AppConfig{
		Meili: MeiliConfig{
			Host:   getEnv("MEILI_HOST", "http://localhost:7700"),
			APIKey: getEnv("MEILI_API_KEY", ""),
		},
		GRPCPort: getEnv("GRPC_PORT", "30005"),
		Brokers:  brokers,
		Topic:    getEnv("KAFKA_TOPIC", "song-events"),
		GroupID:  getEnv("KAFKA_GROUP_ID", "search-service-group"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
