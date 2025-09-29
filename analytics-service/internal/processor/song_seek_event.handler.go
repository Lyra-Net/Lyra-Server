package processor

import (
	analyticsKafka "analytics-service/internal/kafka"
	"context"
	"log"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type SongSeekHandler struct{}

func (h *SongSeekHandler) Handle(ctx context.Context, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) error {
	log.Printf("[SongSeekHandler] Handling message: %s", string(msg))
	return nil
}

func (h *SongSeekHandler) EventType() string {
	return SONG_SEEK_EVENTS
}
