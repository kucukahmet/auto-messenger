package api

import (
	"auto-messager/internal/app"
	"auto-messager/internal/worker"
	"encoding/json"
	"net/http"
)

type Handlers struct {
	app      *app.App
	listener *worker.Listener
}

func NewHandlers(app *app.App, listener *worker.Listener) *Handlers {
	return &Handlers{app: app, listener: listener}
}

func (handler *Handlers) Ping(writer http.ResponseWriter, request *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"message": "pong",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (handler *Handlers) StartListener(writer http.ResponseWriter, request *http.Request) {
	handler.listener.Start()
	response := map[string]interface{}{
		"status":  "ok",
		"message": "Started listener",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (handler *Handlers) StopListener(writer http.ResponseWriter, request *http.Request) {
	handler.listener.Stop()
	response := map[string]interface{}{
		"status":  "ok",
		"message": "Stoped listener",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}
