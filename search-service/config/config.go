package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type MeiliConfig struct {
	Host   string
	APIKey string
}

type AppConfig struct {
	Meili MeiliConfig
	Port  string
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env vars")
	}

	return &AppConfig{
		Meili: MeiliConfig{
			Host:   getEnv("MEILI_HOST", "http://localhost:7700"),
			APIKey: getEnv("MEILI_API_KEY", ""),
		},
		Port: getEnv("PORT", ":3005"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
