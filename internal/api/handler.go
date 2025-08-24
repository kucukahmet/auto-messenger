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

// Ping godoc
// @Summary      Health check
// @Description  Returns simple pong response with server and listener status
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /ping [get]
func (handler *Handlers) Ping(writer http.ResponseWriter, request *http.Request) {

	listenerIsRunning := handler.listener.IsRunning()

	response := map[string]interface{}{
		"status":            "ok",
		"httpServer":        true,
		"listenerIsRunning": listenerIsRunning,
		"message":           "pong",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

// StartListener godoc
// @Summary      Start Listener
// @Description  Starts the automatic message sending loop
// @Produce      json
// @Success      200 {object} map[string]interface{} "started"
// @Failure      409 {string} string "already running"
// @Router       /start-listener [get]
func (handler *Handlers) StartListener(writer http.ResponseWriter, request *http.Request) {
	handler.listener.Start()
	response := map[string]interface{}{
		"status":  "ok",
		"message": "Started listener",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

// StopListener godoc
// @Summary      Stop Listener
// @Description  Stops the automatic message sending loop
// @Produce      json
// @Success      200 {object} map[string]interface{} "stopped"
// @Failure      409 {string} string "not running"
// @Router       /stop-listener [get]
func (handler *Handlers) StopListener(writer http.ResponseWriter, request *http.Request) {
	handler.listener.Stop()
	response := map[string]interface{}{
		"status":  "ok",
		"message": "Stoped listener",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

// ListSentMessages godoc
// @Summary      List sent messages
// @Description  Returns a paginated list of sent messages
// @Tags         messages
// @Produce      json
// @Param        limit  query     int  false "Number of messages to return (default 20, max 100)"
// @Param        offset query     int  false "Offset for pagination (default 0)"
// @Success 200 {array} api.MessageResponse
// @Failure      500    {string}  string "Failed to retrieve messages"
// @Router       /messages/sent [get]
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
