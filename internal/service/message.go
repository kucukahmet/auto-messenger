package service

import (
	"bytes"
	"net/http"
	"time"
)

type MessageService struct {
	EndpointUrl string
	HttpClient  *http.Client
}

func NewMessageService(endpointUrl string) *MessageService {
	return &MessageService{
		EndpointUrl: endpointUrl,
		HttpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (messageService *MessageService) SendMessage(payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, messageService.EndpointUrl, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return messageService.HttpClient.Do(req)
}
