package processor

import (
	analyticsKafka "analytics-service/internal/kafka"
	"context"
	"log"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type SearchHandler struct{}

func (h *SearchHandler) Handle(ctx context.Context, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) error {
	log.Printf("[SearchHandler] Handling message: %s", string(msg))
	return nil
}

func (h *SearchHandler) EventType() string {
	return SEARCH_EVENTS
}
