package api

import (
	"time"
)

type MessageResponse struct {
	ID                int32      `json:"id"`
	PhoneNumber       string     `json:"phone_number"`
	Content           string     `json:"content"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	SentAt            *time.Time `json:"sent_at,omitempty"`
	ResponseMessageID *string    `json:"response_message_id,omitempty"`
	FailReason        *string    `json:"fail_reason,omitempty"`
	RetryCount        int32      `json:"retry_count"`
}

// ErrorResponse swagger:model
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error message"`
}

// SimpleResponse swagger:model
type SimpleResponse struct {
	Status    string `json:"status" example:"ok"`
	Message   string `json:"message" example:"Started listener"`
	IsRunning bool   `json:"isRunning" example:"true"`
}

// HealthResponse swagger:model
type HealthResponse struct {
	Status            string `json:"status" example:"ok"`
	HTTPServer        bool   `json:"httpServer" example:"true"`
	ListenerIsRunning bool   `json:"listenerIsRunning" example:"true"`
	Message           string `json:"message" example:"pong"`
	ApiVersion        string `json:"apiVersion" example:"v1"`
	BuildVersion      string `json:"buildVersion" example:"1.2.3"`
}

// PaginatedMessagesResponse swagger:model
type PaginatedMessagesResponse struct {
	Status        string            `json:"status" example:"ok"`
	Items         []MessageResponse `json:"items"`
	Limit         int               `json:"limit" example:"20"`
	Offset        int               `json:"offset" example:"0"`
	ReturnedCount int               `json:"returned_count" example:"20"`
}

// NewMessageRequest swagger:model
type NewMessageRequest struct {
	PhoneNumber string `json:"phone_number" example:"+905xxxxxxxx"`
	Content     string `json:"content" example:"Hello, this is a test message."`
}

// NewMessageResponse swagger:model
type NewMessageResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Successfully added new message"`
}
