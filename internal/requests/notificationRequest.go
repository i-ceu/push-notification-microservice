package requests

import "push-notification-microservice/internal/models"

type SendNotificationRequest struct {
	PushToken    string          `json:"push_token" binding:"required"`
	UserID       string          `json:"user_id" binding:"required"`
	TemplateCode string          `json:"template_code" binding:"required"`
	ImageURL     string          `json:"image_url,omitempty"`
	ClickURL     string          `json:"click_url,omitempty"`
	Variables    UserData        `json:"variables,omitempty"`
	Metadata     *map[string]any `json:"metadata,omitempty"`
}

type UserData struct {
	Name string          `json:"name,omitempty"`
	Link string          `json:"link,omitempty"`
	Meta *map[string]any `json:"meta,omitempty"`
}

type SendPushNotificationRequest struct {
	PushToken string             `json:"push_token" binding:"required"`
	Title     string             `json:"title" binding:"required"`
	Body      string             `json:"body" binding:"required"`
	ImageURL  string             `json:"image_url,omitempty"`
	ClickURL  string             `json:"click_url,omitempty"`
	Data      models.MessageData `json:"data,omitempty"`
}
