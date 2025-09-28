package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func NewConsumer(brokers, groupID, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokers},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1e3,
		MaxBytes: 10e6,
	})
}

func ConsumeLoop(ctx context.Context, reader *kafka.Reader, handler func([]byte)) {
	for {
		select {
		case <-ctx.Done():
			log.Println("ConsumeLoop stopped (context canceled)")
			return
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("ConsumeLoop exiting due to context cancellation")
					return
				}
				log.Printf("consumer error: %v", err)
				continue
			}
			handler(m.Value)
		}
	}
}
