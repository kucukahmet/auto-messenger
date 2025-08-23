package api

import (
	"auto-messager/internal/worker"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(listener *worker.Listener, db *sql.DB) http.Handler {
	router := chi.NewRouter()
	handler := NewHandlers(listener, db)

	router.Get("/ping", handler.Ping)
	router.Get("/start-listener", handler.StartListener)
	router.Get("/stop-listener", handler.StopListener)

	return router
}
