package models

import "time"

type PushNotification struct {
	CorrelationID   string          `json:"correlation_id"`
	NotificationID  string          `json:"notification_id"`
	UserData        map[string]any  `json:"user_id"`
	PushToken       string          `json:"push_token"`
	RenderedContent RenderedContent `json:"data,omitempty"`
	RetryCount      int             `json:"retry_count"`
	CreatedAt       time.Time       `json:"created_at"`
}

type RenderedContent struct {
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Variables map[string]any `json:"variables"`
}
