package clickhouse

import (
	"context"
	"fmt"
	"log"

	"analytics-service/config"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

func New(cfg *config.Config) ch.Conn {
	addr := fmt.Sprintf("%s:%s", cfg.ClickHouseHost, cfg.ClickHousePort)

	conn, err := ch.Open(&ch.Options{
		Addr: []string{addr},
		Auth: ch.Auth{
			Database: cfg.ClickHouseDB,
			Username: cfg.ClickHouseUser,
			Password: cfg.ClickHousePass,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect ClickHouse: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("clickhouse ping error: %v", err)
	}

	log.Println("Connected to ClickHouse")
	return conn
}
