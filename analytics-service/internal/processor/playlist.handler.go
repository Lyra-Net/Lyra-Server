package processor

import (
	analyticsKafka "analytics-service/internal/kafka"
	"context"
	"log"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type PlaylistHandler struct{}

func (h *PlaylistHandler) Handle(ctx context.Context, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) error {
	log.Printf("[PlaylistHandler] Handling message: %s", string(msg))
	return nil
}

func (h *PlaylistHandler) EventType() string {
	return PLAYLIST_EVENTS
}
