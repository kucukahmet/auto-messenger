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
// @Description  Returns simple pong response with server and listener status. You can also learn the API version.
// @Tags         health
// @Produce      json
// @Success      200 {object} api.HealthResponse
// @Router       /ping [get]
func (handler *Handlers) Ping(writer http.ResponseWriter, request *http.Request) {
	listenerIsRunning := handler.listener.IsRunning()

	response := map[string]interface{}{
		"status":            "ok",
		"httpServer":        true,
		"listenerIsRunning": listenerIsRunning,
		"message":           "pong",
		"apiVersion":        handler.app.Config.API_VERSION,
		"buildVersion":      handler.app.Config.BUILD_VERSION,
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

// StartListener godoc
// @Summary      Start Listener
// @Description  Starts the automatic message sending loop
// @Tags         listener
// @Produce      json
// @Success      200 {object} api.SimpleResponse "started"
// @Failure      409 {object} api.ErrorResponse "already running"
// @Router       /api/start-listener [get]
func (handler *Handlers) StartListener(writer http.ResponseWriter, request *http.Request) {
	if handler.listener.IsRunning() {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(writer).Encode(ErrorResponse{
			Status:  "error",
			Message: "Already running",
		})
		return
	}

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
// @Tags         listener
// @Produce      json
// @Success      200 {object} api.SimpleResponse "stopped"
// @Failure      409 {object} api.ErrorResponse "not running"
// @Router       /api/stop-listener [get]
func (handler *Handlers) StopListener(writer http.ResponseWriter, request *http.Request) {
	if !handler.listener.IsRunning() {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(writer).Encode(ErrorResponse{
			Status:  "error",
			Message: "Not running",
		})
		return
	}

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
// @Param        limit   query     int  false "Number of messages to return (default 20, max 100)"
// @Param        offset  query     int  false "Offset for pagination (default 0)"
// @Success      200     {object}  api.PaginatedMessagesResponse
// @Failure      500     {object}  api.ErrorResponse "Failed to retrieve messages"
// @Router       /api/messages/sent [get]
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

// AddNewMessage godoc
// @Summary      Add new message
// @Description  Adds a new message to the queue for sending
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        message  body      api.NewMessageRequest	true  "Message to add"
// @Success      201      {object}  api.NewMessageResponse "Message added to queue"
// @Failure      400      {object}  api.ErrorResponse    "Invalid request payload"
// @Failure      500      {object}  api.ErrorResponse    "Failed to add
// @Router       /api/messages [post]
func (handler *Handlers) AddNewMessage(writer http.ResponseWriter, request *http.Request) {
	var newMessageReq NewMessageRequest
	err := json.NewDecoder(request.Body).Decode(&newMessageReq)
	if err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = handler.app.Queries.InsertMessage(request.Context(), storage.InsertMessageParams{
		PhoneNumber: newMessageReq.PhoneNumber,
		Content:     newMessageReq.Content,
	})

	if err != nil {
		http.Error(writer, "Failed to add message to queue", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Message added to queue",
	})
}
