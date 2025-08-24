package api

import (
	"auto-messager/internal/app"
	"auto-messager/internal/worker"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "auto-messager/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(app *app.App, listener *worker.Listener) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	handler := NewHandlers(app, listener)

	router.Get("/ping", handler.Ping)
	router.Get("/start-listener", handler.StartListener)
	router.Get("/stop-listener", handler.StopListener)
	router.Get("/messages/sent", handler.ListSentMessages)

	// Swagger documentation
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	return router
}
