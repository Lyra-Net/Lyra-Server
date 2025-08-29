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

	go func() {
		log.Println("Http server started on ", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Http server error: %v", err)
		}
	}()

	<-s.ctx.Done()
	log.Println("Shutting down http server...")

	return s.server.Shutdown(context.Background())
}
