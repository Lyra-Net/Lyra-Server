package processor

import (
	analyticsKafka "analytics-service/internal/kafka"
	"context"
	"log"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type SongPlayHandler struct{}

func (h *SongPlayHandler) Handle(ctx context.Context, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) error {
	log.Printf("[SongPlayHandler] Handling message: %s", string(msg))
	return nil
}

func (h *SongPlayHandler) EventType() string {
	return SONG_PLAY_EVENTS
}
