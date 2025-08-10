package routers

import "github.com/go-chi/chi/v5"

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	streamRouter := NewStreamRouter()
	streamRouter.Init(r)
	return r
}
