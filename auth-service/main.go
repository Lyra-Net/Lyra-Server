package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	"auth-service/config"
	"auth-service/internal/interceptor"
	"auth-service/internal/repository"
	"auth-service/utils"

	"github.com/trandinh0506/BypassBeats/proto/gen/auth"

	"auth-service/internal/redisconn"
	"auth-service/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.GetConfig()
	utils.InitCrypto()

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
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor()),
	)

	auth.RegisterAuthServiceServer(grpcServer, services.NewAuthServer(queries))
	reflection.Register(grpcServer)
	fmt.Println("gRPC server is running on port " + cfg.PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
