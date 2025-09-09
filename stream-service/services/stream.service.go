package services

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type StreamService struct {
	minio *MinioService
}

func NewStreamService(minio *MinioService) *StreamService {
	return &StreamService{minio: minio}
}

func (s *StreamService) StreamFile(w http.ResponseWriter, r *http.Request, filename string) error {
	ctx := context.Background()

	object, err := s.minio.Client.GetObject(ctx, s.minio.Bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		log.Println("error getting object:", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return err
	}
	defer object.Close()

	w.Header().Set("Content-Type", "audio/mpeg")

	if _, err := io.Copy(w, object); err != nil {
		log.Println("error writing response:", err)
		return err
	}

	return nil
}
