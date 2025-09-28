package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"analytics-service/config"
	"analytics-service/internal/clickhouse"
	"analytics-service/internal/kafka"
	"analytics-service/internal/processor"
)

func main() {
	cfg := config.GetConfig()

	chConn := clickhouse.New(cfg)

	log.Println("Starting Kafka consumer...")
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.GroupID, cfg.InputTopic)
	defer consumer.Close()

	log.Println("Starting Kafka producer...")
	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.OutputTopic)
	defer producer.Close()

	log.Println("Analytics Service started...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v. Shutting down gracefully...\n", sig)
		cancel()
	}()

	kafka.ConsumeLoop(ctx, consumer, func(msg []byte) {
		processor.HandleEvent(chConn, producer, msg)
	})

	log.Println("Service stopped.")
}
