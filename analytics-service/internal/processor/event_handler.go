package processor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"analytics-service/internal/kafka"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type PlayEvent struct {
	UserID      string    `json:"user_id"`
	SongID      string    `json:"song_id"`
	DeviceID    string    `json:"device_id"`
	Timestamp   time.Time `json:"timestamp"`
	StartSecond int       `json:"start_second"`
	EndSecond   int       `json:"end_second"`
}

func HandleEvent(conn ch.Conn, producer *kafka.Producer, raw []byte) {
	var evt PlayEvent
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("invalid event: %v", err)
		return
	}

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO song_play_events (user_id, song_id, device_id, ts, duration)")
	if err != nil {
		log.Printf("prepare batch error: %v", err)
		return
	}

	if err := batch.Append(evt.UserID, evt.SongID, evt.DeviceID, evt.Timestamp, evt.StartSecond, evt.EndSecond); err != nil {
		log.Printf("append error: %v", err)
		return
	}

	if err := batch.Send(); err != nil {
		log.Printf("send error: %v", err)
		return
	}

	log.Printf("Inserted event for song %s", evt.SongID)

	update := map[string]interface{}{
		"song_id": evt.SongID,
		"inc":     1,
	}
	b, _ := json.Marshal(update)
	producer.Publish(context.Background(), b)
}
