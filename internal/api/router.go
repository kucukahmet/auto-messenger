package api

import (
	"auto-messager/internal/app"
	"auto-messager/internal/worker"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(app *app.App, listener *worker.Listener) http.Handler {
	router := chi.NewRouter()
	handler := NewHandlers(app, listener)

	router.Get("/ping", handler.Ping)
	router.Get("/start-listener", handler.StartListener)
	router.Get("/stop-listener", handler.StopListener)

	return router
}
