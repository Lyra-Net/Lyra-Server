package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	ctx    context.Context
	server *http.Server
}

func NewHttpServer(ctx context.Context, addr string, router *gin.Engine) *HttpServer {
	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: router,
	}
	return &HttpServer{
		ctx:    ctx,
		server: srv,
	}
}

func (s *HttpServer) Start() error {
	log.Println("HTTP server started on", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}
