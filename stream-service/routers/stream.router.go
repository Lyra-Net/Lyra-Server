package routers

import (
	"stream-service/handlers"

	"github.com/go-chi/chi/v5"
)

type StreamRouter struct {
	streamHandler *handlers.StreamHandler
}

func NewStreamRouter() *StreamRouter {
	return &StreamRouter{
		streamHandler: handlers.NewStreamHandler(),
	}
}

func (streamRouter *StreamRouter) Init(r *chi.Mux) {
	r.Get("/stream/{filename}", streamRouter.streamHandler.StreamAudio)
}
