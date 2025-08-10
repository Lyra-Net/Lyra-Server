package services

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type StreamService struct {
}

func NewStreamService() *StreamService {
	return &StreamService{}
}

func (s *StreamService) StreamFile(w http.ResponseWriter, r *http.Request, filename string) error {
	filePath := filepath.Join("audios", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Type", "audio/mp3")
	http.ServeContent(w, r, filename, fileStat(filePath), file)
	return nil
}

func fileStat(path string) (modTime time.Time) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	return info.ModTime()
}
