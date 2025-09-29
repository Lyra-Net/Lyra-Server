package processor

import (
	"context"
	"log"

	analyticsKafka "analytics-service/internal/kafka"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

var handlers = map[string]EventHandler{}

func RegisterHandler(handler EventHandler) {
	handlers[handler.EventType()] = handler
}

func HandleEvent(ctx context.Context, topic string, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) {
	h, ok := handlers[topic]
	if !ok {
		log.Printf("No handler registered for topic %s", topic)
		return
	}
	if err := h.Handle(ctx, msg, chConn, producer); err != nil {
		log.Printf("Error handling message for topic %s: %v", topic, err)
	}
}
