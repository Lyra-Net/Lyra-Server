package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers string
	InputTopics  []string
	OutputTopic  string
	GroupID      string

	ClickHouseDSN  string
	ClickHouseHost string
	ClickHousePort string
	ClickHouseUser string
	ClickHousePass string
	ClickHouseDB   string
}

var (
	appEnv *Config
	once   sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, using system env only")
		}

		KafkaBrokers := getEnv("KAFKA_BROKERS", "kafka:9092")
		rawTopics := getEnv("KAFKA_INPUT_TOPICS", "song_play_events")
		InputTopics := parseTopics(rawTopics)
		OutputTopic := getEnv("KAFKA_OUTPUT_TOPIC", "song_play_counts")
		GroupID := getEnv("KAFKA_GROUP_ID", "analytics-consumer")

		ClickHouseDSN := getEnv("CLICKHOUSE_DSN", "")
		ClickHouseHost := getEnv("CLICKHOUSE_HOST", "clickhouse")
		ClickHousePort := getEnv("CLICKHOUSE_PORT", "9000")
		ClickHouseUser := getEnv("CLICKHOUSE_USER", "default")
		ClickHousePass := getEnv("CLICKHOUSE_PASSWORD", "")
		ClickHouseDB := getEnv("CLICKHOUSE_DB", "default")

		if KafkaBrokers == "" || ClickHouseDSN == "" || len(InputTopics) == 0 || OutputTopic == "" || GroupID == "" {
			log.Fatal("Some required env vars are missing")
		}

		appEnv = &Config{
			KafkaBrokers:   KafkaBrokers,
			ClickHouseDSN:  ClickHouseDSN,
			InputTopics:    InputTopics,
			OutputTopic:    OutputTopic,
			GroupID:        GroupID,
			ClickHouseHost: ClickHouseHost,
			ClickHousePort: ClickHousePort,
			ClickHouseUser: ClickHouseUser,
			ClickHousePass: ClickHousePass,
			ClickHouseDB:   ClickHouseDB,
		}
	})

	return appEnv
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseTopics(raw string) []string {
	topics := strings.Split(raw, ",")
	for i := range topics {
		topics[i] = strings.TrimSpace(topics[i])
	}
	return topics
}
