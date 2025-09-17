package mq

import (
	"context"
	"encoding/json"
	"log"
	"search-service/dto"
	"search-service/server"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader      *kafka.Reader
	meiliClient *server.MeiliClient
}

func NewKafkaConsumer(brokers []string, topic, groupID string, meiliClient *server.MeiliClient) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		meiliClient: meiliClient,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading kafka message: %v", err)
			continue
		}

		var event struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}

		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		switch event.Type {
		case "song_created":
			var song dto.CreateSongRequest
			if err := json.Unmarshal(event.Payload, &song); err != nil {
				log.Printf("Failed to unmarshal song: %v", err)
				continue
			}
			log.Printf("New song created: %+v", song)

			if err := c.meiliClient.IndexSong(song); err != nil {
				log.Printf("Failed to index song: %v", err)
			}

		case "artist_created":
			var artist dto.Artist
			if err := json.Unmarshal(event.Payload, &artist); err != nil {
				log.Printf("Failed to unmarshal artist: %v", err)
				continue
			}
			log.Printf("New artist created: %+v", artist)

			if err := c.meiliClient.IndexArtist(artist); err != nil {
				log.Printf("Failed to index artist: %v", err)
			}

		default:
			log.Printf("Unknown event type: %s", event.Type)
		}
	}
}
