package api

import (
	"auto-messager/internal/worker"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(listener *worker.Listener) http.Handler {
	router := chi.NewRouter()
	handler := NewHandlers(listener)

	router.Get("/ping", handler.Ping)
	router.Get("/start-listener", handler.StartListener)

	return router
}
