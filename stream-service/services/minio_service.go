package services

import (
	"log"
	"stream-service/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	Client *minio.Client
	Bucket string
}

func NewMinioService() *MinioService {
	cfg := config.GetConfig()
	useSSL := false

	log.Println("MinIO Config:", cfg.MINIO_ENDPOINT, cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY, cfg.MINIO_BUCKET)
	//
	credentials := credentials.NewStaticV4(cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY, "")

	client, err := minio.New(cfg.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials,
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to init MinIO: %v", err)
	}

	return &MinioService{
		Client: client,
		Bucket: cfg.MINIO_BUCKET,
	}
}
