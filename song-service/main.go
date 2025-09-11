package main

import (
	"context"
	"log"
	"os/signal"
	"song-service/config"
	mq "song-service/internal/MQ"
	"song-service/internal/playlist"
	"song-service/internal/repository"
	"song-service/internal/router"
	"song-service/internal/server"
	"syscall"
	"time"

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

	// Error channel
	errCh := make(chan error, 2)

	// Kafka producer
	log.Println("starting kafka producer...")
	producer := mq.NewKafkaProducer(cfg.Brokers, cfg.Topic)
	defer func() {
		log.Println("closing kafka producer...")
		if err := producer.Close(); err != nil {
			log.Printf("Error when closing kafka producer: %v", err)
		}
		log.Println("kafka producer closed")
	}()

	// HTTP server
	log.Println("starting http server...")
	r := router.NewRouter(q, producer)
	httpServer := server.NewHttpServer(ctx, cfg.HttpPort, r)

	go func() {
		if err := httpServer.Start(); err != nil {
			errCh <- err
		}
	}()

	// gRPC server
	log.Println("starting gRPC server...")
	playlistSvc := playlist.NewPlaylistService(q)
	grpcServer := server.NewGrpcServer(cfg.GRPCPort, playlistSvc)

	go func() {
		if err := grpcServer.Start(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down song-service...")
		grpcServer.Stop()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Stop(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("song-service stopped")
	case err := <-errCh:
		log.Printf("song-serice error: %v", err)
	}
}
