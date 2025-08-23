package api

import (
	"auto-messager/internal/worker"
	"encoding/json"
	"net/http"
)

type Handlers struct {
	listener *worker.Listener
}

func NewHandlers(listener *worker.Listener) *Handlers {
	return &Handlers{listener: listener}
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
