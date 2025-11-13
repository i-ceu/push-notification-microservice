package models

import "time"

type PushNotification struct {
	CorrelationID  string      `json:"correlation_id"`
	NotificationID string      `json:"notification_id"`
	UserID         string      `json:"user_id"`
	PushToken      string      `json:"push_token"`
	Data           MessageData `json:"data,omitempty"`
	RetryCount     int         `json:"retry_count"`
	CreatedAt      time.Time   `json:"created_at"`
}

type MessageData struct {
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Variables map[string]any `json:"variables"`
}
