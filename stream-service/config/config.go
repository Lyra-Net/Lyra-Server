package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Cfg struct {
	PORT             string
	MINIO_ENDPOINT   string
	MINIO_ACCESS_KEY string
	MINIO_SECRET_KEY string
	MINIO_BUCKET     string
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
		minioEndpoint := os.Getenv("MINIO_ENDPOINT")
		minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
		minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
		minioBucket := os.Getenv("MINIO_BUCKET")

		if port == "" || minioEndpoint == "" || minioAccessKey == "" || minioSecretKey == "" || minioBucket == "" {
			log.Fatal("Some required env vars are missing")
		}

		appEnv = &Cfg{
			PORT:             port,
			MINIO_ENDPOINT:   minioEndpoint,
			MINIO_ACCESS_KEY: minioAccessKey,
			MINIO_SECRET_KEY: minioSecretKey,
			MINIO_BUCKET:     minioBucket,
		}
	})

	return appEnv
}
