package services

import (
	"song-service/models"
	"song-service/repository"

	"gorm.io/gorm"
)

func CreateArtist(db *gorm.DB, artist *models.Artist) error {
	return repository.CreateArtist(db, artist)
}

func GetArtists(db *gorm.DB, page, limit int, name string) ([]models.Artist, int64, error) {
	return repository.GetArtists(db, page, limit, name)
}
