package main

import (
	"log"
	"net"

	"search-service/config"
	"search-service/proto/search"
	"search-service/server"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	search.RegisterSearchServiceServer(s, server.NewSearchServer(cfg.Meili))

	log.Printf("gRPC server is running on %s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
