package api

import (
	"auto-messager/internal/app"
	"auto-messager/internal/storage"
	"auto-messager/internal/worker"
	"encoding/json"
	"net/http"
	"strconv"
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

func (handler *Handlers) ListSentMessages(writer http.ResponseWriter, request *http.Request) {
	limitStr := request.URL.Query().Get("limit")
	offsetStr := request.URL.Query().Get("offset")

	const defaultLimit = 20
	const defaultOffset = 0
	const maxLimit = 100

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	messages, err := handler.app.Queries.ListSent(request.Context(), storage.ListSentParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		http.Error(writer, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(messages)
}
