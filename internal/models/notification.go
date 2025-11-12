package models

import "time"

type PushNotification struct {
	CorrelationID  string         `json:"correlation_id"`
	NotificationID string         `json:"notification_id"`
	UserID         string         `json:"user_id"`
	PushToken      string         `json:"push_token"`
	TemplateCode   string         `json:"template_code"`
	Title          string         `json:"title"`
	Body           string         `json:"body"`
	ImageURL       string         `json:"image_url,omitempty"`
	ClickURL       string         `json:"click_url,omitempty"`
	Data           map[string]any `json:"data,omitempty"`
	Variables      map[string]any `json:"variables"`
	RetryCount     int            `json:"retry_count"`
	CreatedAt      time.Time      `json:"created_at"`
}
