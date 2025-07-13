package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"identity-service/config"
	"identity-service/models"
	"identity-service/proto/auth"
	"identity-service/redisconn"
	"identity-service/services"

	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()
	config.DB.AutoMigrate(&models.User{}, &models.RefreshToken{})
	redisconn.InitRedis()
	PORT := os.Getenv("PORT")
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, services.NewAuthServer(config.DB))

	fmt.Println("gRPC server is running on port " + PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
