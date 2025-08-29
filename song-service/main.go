package main

import (
	"context"
	"log"
	"os/signal"
	"song-service/config"
	mq "song-service/internal/MQ"
	"song-service/internal/repository"
	"song-service/internal/router"
	"song-service/internal/server"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadConfig()

	pool, err := pgxpool.New(ctx, cfg.DbURL)

	if err != nil {
		log.Fatalf("Error when opening DB. err: %v", err)
	}

	defer pool.Close()

	q := repository.New(pool)

	errCh := make(chan error, 2)

	producer := mq.NewKafkaProducer(cfg.Brokers, cfg.Topic)

	r := router.NewRouter(q, producer)

	httpServer := server.NewHttpServer(ctx, cfg.HttpPort, r)

	go func() {
		if err := httpServer.Start(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown song-service")
	case err := <-errCh:
		log.Printf("song-serice error: %v", err)
	}
}
