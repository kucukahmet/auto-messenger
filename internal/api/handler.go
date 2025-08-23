package api

import (
	"encoding/json"
	"net/http"
)

func Ping(writer http.ResponseWriter, request *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"message": "pong",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}
