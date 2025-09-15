package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"identity-service/config"
	"identity-service/internal/repository"
	"identity-service/proto/auth"
	"identity-service/redisconn"
	"identity-service/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.GetConfig()

	// Redis
	redisconn.InitRedis()

	// Postgres
	pool, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	queries := repository.New(pool)

	// gRPC
	lis, err := net.Listen("tcp", ":"+cfg.PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, services.NewAuthServer(queries))

	fmt.Println("gRPC server is running on port " + cfg.PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
