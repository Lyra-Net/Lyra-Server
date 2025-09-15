package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Cfg struct {
	PORT               string
	DB_URL             string
	POSTGRES_HOST      string
	POSTGRES_USER      string
	POSTGRES_PASSWORD  string
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
	REDIS_URL          string
	REDIS_HOST         string
	REDIS_PORT         string
}

var (
	appEnv *Cfg
	once   sync.Once
)

func GetConfig() *Cfg {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, using system env only")
		}

		port := os.Getenv("PORT")
		dbURL := os.Getenv("DB_URL")
		postgresHost := os.Getenv("POSTGRES_HOST")
		postgresUser := os.Getenv("POSTGRES_USER")
		postgresPassword := os.Getenv("POSTGRES_PASSWORD")
		jwtAccess := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
		jwtRefresh := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
		redisURL := os.Getenv("REDIS_URL")
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")

		if port == "" || dbURL == "" || jwtAccess == "" || jwtRefresh == "" || redisURL == "" {
			log.Fatal("Some required env vars are missing")
		}

		appEnv = &Cfg{
			PORT:               port,
			DB_URL:             dbURL,
			POSTGRES_HOST:      postgresHost,
			POSTGRES_USER:      postgresUser,
			POSTGRES_PASSWORD:  postgresPassword,
			JWT_ACCESS_SECRET:  jwtAccess,
			JWT_REFRESH_SECRET: jwtRefresh,
			REDIS_URL:          redisURL,
			REDIS_HOST:         redisHost,
			REDIS_PORT:         redisPort,
		}
	})

	return appEnv
}
