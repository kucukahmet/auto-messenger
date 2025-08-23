package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() http.Handler {
	router := chi.NewRouter()

	router.Get("/ping", Ping)

	return router
}
