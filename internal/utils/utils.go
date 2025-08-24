package utils

import (
	"auto-messager/internal/storage"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OutboundPayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func BuildPayloadFromMessage(msg *storage.Message) ([]byte, error) {
	p := OutboundPayload{
		To:      msg.PhoneNumber,
		Content: msg.Content,
	}
	return json.Marshal(p)
}

func ParseWebhookResponse(resp *http.Response) (*WebhookResponse, error) {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("webhook returned %d: %s", resp.StatusCode, string(body))
	}
	var wr WebhookResponse
	if err := json.Unmarshal(body, &wr); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w; raw=%s", err, string(body))
	}
	return &wr, nil
}
