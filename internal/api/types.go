package api

import "time"

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
