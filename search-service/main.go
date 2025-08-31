package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"search-service/config"
	"search-service/mq"
	"search-service/server"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	// Init meili client
	meiliClient := server.NewMeiliClient(cfg.Meili)

	// Init gRPC server
	go func() {
		if err := server.RunGRPCServer(":"+cfg.GRPCPort, meiliClient); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()
	log.Println("gRPC server started on :", cfg.GRPCPort)

	// Init Kafka consumer
	log.Println("Starting Kafka consumer....")
	kafkaConsumer := mq.NewKafkaConsumer(cfg.Brokers, cfg.Topic, cfg.GroupID, meiliClient)

	ctx, cancel := context.WithCancel(context.Background())
	go kafkaConsumer.Start(ctx)
	log.Println("Kafka consumer started")

	// Wait for termination signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	cancel()
}
