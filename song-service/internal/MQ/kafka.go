package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Emit(ctx context.Context, eventType string, payload interface{}) error {
	data, err := json.Marshal(struct {
		Type      string      `json:"type"`
		Payload   interface{} `json:"payload"`
		Timestamp time.Time   `json:"timestamp"`
	}{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(eventType),
		Value: data,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
