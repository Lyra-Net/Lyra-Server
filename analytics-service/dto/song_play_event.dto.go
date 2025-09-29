package dto

import (
	"time"

	"github.com/google/uuid"
)

type SongPlayEvent struct {
	UserID      uuid.UUID `json:"user_id"`
	SongID      string    `json:"song_id"`
	DeviceID    uuid.UUID `json:"device_id"`
	Timestamp   time.Time `json:"timestamp"`
	StartSecond int       `json:"start_second"`
	EndSecond   int       `json:"end_second"`
}
