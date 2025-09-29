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

	log.Println("Starting Kafka consumers...")
	consumers := kafka.NewConsumers(cfg.KafkaBrokers, cfg.GroupID, cfg.InputTopics)

	log.Println("Starting Kafka producer...")
	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.OutputTopic)

	processor.RegisterHandler(&processor.SongPlayHandler{})
	processor.RegisterHandler(&processor.PlaylistHandler{})
	processor.RegisterHandler(&processor.SearchHandler{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v. Shutting down...\n", sig)
		cancel()
	}()

	kafka.ConsumeMultipleTopics(ctx, consumers, func(msg []byte, topic string) {
		processor.HandleEvent(ctx, topic, msg, chConn, producer)
	})

	log.Println("Stopping consumers and producer...")
	kafka.CloseConsumers(consumers)
	producer.Close()
	log.Println("Analytics Service stopped.")
}
