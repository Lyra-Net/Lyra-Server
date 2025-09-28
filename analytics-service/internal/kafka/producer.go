package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(ctx context.Context, value []byte) error {
	err := p.writer.WriteMessages(ctx,
		kafka.Message{Value: value},
	)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("publish canceled due to context cancellation")
			return ctx.Err()
		}
		log.Printf("publish error: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	log.Println("Closing Kafka producer...")
	return p.writer.Close()
}
