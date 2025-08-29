package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DbURL    string
	HttpPort string
	GRPCPort string
	Brokers  []string
	Topic    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	brokers := strings.Split(os.Getenv("KAFKA_BROKER"), ",")
	return &Config{
		DbURL:    os.Getenv("DB_URL"),
		HttpPort: os.Getenv("HTTP_PORT"),
		GRPCPort: os.Getenv("GRPC_PORT"),
		Brokers:  brokers,
		Topic:    os.Getenv("KAFKA_TOPIC"),
	}
}
