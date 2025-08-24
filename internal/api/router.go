package api

import (
	"auto-messager/internal/app"
	"auto-messager/internal/worker"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "auto-messager/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func BuildVersionMiddleware(buildVersion string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if buildVersion != "" {
				w.Header().Set("X-Build-Version", buildVersion)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Router(app *app.App, listener *worker.Listener) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(BuildVersionMiddleware(app.Config.BUILD_VERSION))

	handler := NewHandlers(app, listener)

	router.Get("/ping", handler.Ping)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/api", func(api chi.Router) {
		api.Get("/start-listener", handler.StartListener)
		api.Get("/stop-listener", handler.StopListener)
		api.Get("/messages/sent", handler.ListSentMessages)
		api.Post("/messages", handler.AddNewMessage)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":{"code":"not_found","message":"route not found"}}`)
	})

	return router
}
