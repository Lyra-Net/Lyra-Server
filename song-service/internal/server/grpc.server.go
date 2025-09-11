package server

import (
	"log"
	"net"

	pb "github.com/trandinh0506/BypassBeats/proto/gen/playlist"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	port   string
	server *grpc.Server
}

func NewGrpcServer(port string, svc pb.PlaylistServiceServer) *GrpcServer {
	grpcServer := grpc.NewServer()
	pb.RegisterPlaylistServiceServer(grpcServer, svc)
	reflection.Register(grpcServer)

	return &GrpcServer{
		port:   port,
		server: grpcServer,
	}
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	log.Printf("gRPC server started at %s", s.port)
	return s.server.Serve(lis)
}

func (s *GrpcServer) Stop() {
	log.Println("Stopping gRPC server...")
	s.server.GracefulStop()
}
