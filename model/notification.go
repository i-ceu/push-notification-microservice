package model

import "time"

type PushNotification struct {
	CorrelationID      string    `json:"correlation_id"`
	NotificationID     string    `json:"notification_id"`
	UserData           string    `json:"user_id"`
	PushToken          string    `json:"push_token"`
	NotificattionTitle string    `json:"notification_title,omitempty"`
	NotificattionBody  string    `json:"notification_body,omitempty"`
	RetryCount         int       `json:"retry_count"`
	CreatedAt          time.Time `json:"created_at"`
}

// type RenderedContent struct {
// 	Title     string         `json:"title"`
// 	Body      string         `json:"body"`
// 	Variables map[string]any `json:"variables"`
// }

// type TemplateData struct {
// 	Template_code string `json:"template_code"`
// 	Name          string `json:"name"`
// 	Link          string `json:"link"`
// }
// type QueueResponse struct {
// 	Correlation_id string       `json:"correlation_id"` //for tracking
// 	Data           TemplateData `json:"data"`
// 	PushToken      string       `json:"push_token"`
// 	Title          string       `json:"title"`
// 	Body           string       `json:"body"`
// }
