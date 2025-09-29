package kafka

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

func NewConsumers(brokers, groupID string, topics []string) []*kafka.Reader {
	readers := make([]*kafka.Reader, 0, len(topics))
	for _, topic := range topics {
		readers = append(readers, kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{brokers},
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 1e3,
			MaxBytes: 10e6,
		}))

	}
	return readers
}

func ConsumeMultipleTopics(ctx context.Context, readers []*kafka.Reader, handler func([]byte, string)) {
	var wg sync.WaitGroup
	wg.Add(len(readers))

	for _, r := range readers {
		go func(r *kafka.Reader) {
			topic := r.Config().Topic
			defer wg.Done()
			log.Printf("Consumer started for topic: %s\n", topic)
			for {
				select {
				case <-ctx.Done():
					log.Printf("Consumer for topic %s stopping (context canceled)\n", topic)
					return
				default:
					msg, err := r.ReadMessage(ctx)
					if err != nil {
						if ctx.Err() != nil {
							log.Printf("Consumer for topic %s exiting due to context cancellation\n", topic)
							return
						}
						log.Printf("Consumer error on topic %s: %v", topic, err)
						continue
					}
					handler(msg.Value, topic)
				}
			}
		}(r)
	}

	wg.Wait()
}

func CloseConsumers(readers []*kafka.Reader) {

	for _, r := range readers {
		topic := r.Config().Topic
		if err := r.Close(); err != nil {
			log.Printf("Error closing consumer for topic %s: %v", topic, err)
		} else {
			log.Printf("Consumer closed for topic %s", topic)
		}
	}
}
