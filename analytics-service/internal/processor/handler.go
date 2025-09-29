package processor

import (
	analyticsKafka "analytics-service/internal/kafka"
	"context"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

const (
	SONG_PLAY_EVENTS = "song_play_events"
	PLAYLIST_EVENTS  = "playlist_events"
	SEARCH_EVENTS    = "search_events"
)

type EventHandler interface {
	Handle(ctx context.Context, msg []byte, chConn ch.Conn, producer *analyticsKafka.Producer) error
	EventType() string
}
