package config

import (
	"log"
	"os"
	"song-service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB          *gorm.DB
	KafkaBroker string
	Port        string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	// 👇 Tự động tạo bảng
	if err := db.AutoMigrate(&models.Song{}, &models.Playlist{}, &models.PlaylistItem{}); err != nil {
		log.Fatal("Failed to migrate DB:", err)
	}

	return &Config{
		DB:          db,
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		Port:        os.Getenv("PORT"),
	}
}
