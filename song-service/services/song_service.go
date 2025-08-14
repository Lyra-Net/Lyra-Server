package services

import (
	"song-service/models"
	"song-service/repository"

	"gorm.io/gorm"
)

func CreateSong(db *gorm.DB, song *models.Song) error {
	return repository.SaveSong(db, song)
}

func GetSongs(db *gorm.DB, page, limit int) ([]models.Song, int64, error) {
	return repository.GetSongs(db, page, limit)
}

func GetSongById(db *gorm.DB, id string) (models.Song, error) {
	return repository.GetSongById(db, id)
}
