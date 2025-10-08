package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Mailer struct {
	From     string
	Host     string
	Port     string
	Password string
}

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
	MAILER             *Mailer
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
		mailFrom := os.Getenv("SMTP_EMAIL")
		mailHost := os.Getenv("SMTP_HOST")
		mailPort := os.Getenv("SMTP_PORT")
		mailPass := os.Getenv("SMTP_PASSWORD")
		mailer := &Mailer{
			From:     mailFrom,
			Host:     mailHost,
			Port:     mailPort,
			Password: mailPass,
		}
		if port == "" || dbURL == "" ||
			jwtAccess == "" || jwtRefresh == "" ||
			redisURL == "" || mailFrom == "" ||
			mailHost == "" || mailPort == "" || mailPass == "" {
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
			MAILER:             mailer,
		}
	})

	return appEnv
}
