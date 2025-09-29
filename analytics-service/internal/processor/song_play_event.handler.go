package processor

import (
	"analytics-service/dto"
	analyticsKafka "analytics-service/internal/kafka"
	"context"
	"encoding/json"
	"log"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type SongPlayHandler struct{}

func (h *SongPlayHandler) Handle(ctx context.Context, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) error {
	log.Printf("[SongPlayHandler] Handling message: %s", string(msg))

	var req dto.SongPlayEvent

	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println("error when unmarshal message: ", err)
		return err
	}

	log.Println("req data: ", req)

	return nil
}

func (h *SongPlayHandler) EventType() string {
	return SONG_PLAY_EVENTS
}
