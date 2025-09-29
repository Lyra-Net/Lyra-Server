package dto

import (
	"time"

	"github.com/google/uuid"
)

type SongSeekEvent struct {
	UserID     uuid.UUID `json:"user_id"`
	SongID     string    `json:"song_id"`
	DeviceID   uuid.UUID `json:"device_id"`
	FromSecond int       `json:"from_second"`
	ToSecond   int       `json:"to_second"`
	Timestamp  time.Time `json:"timestamp"`
}
