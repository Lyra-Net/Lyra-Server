package handlers

import (
	"net/http"
	"stream-service/services"

	"github.com/go-chi/chi/v5"
)

type StreamHandler struct {
	streamService *services.StreamService
}

func NewStreamHandler() *StreamHandler {
	minIOService := services.NewMinioService()
	return &StreamHandler{
		streamService: services.NewStreamService(minIOService),
	}
}

func (h *StreamHandler) StreamAudio(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	err := h.streamService.StreamFile(w, r, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
