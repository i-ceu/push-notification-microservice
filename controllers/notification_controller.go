package controllers

import (
	"net/http"
	"push_service/model"
	"push_service/push"
	"push_service/requests"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PushController struct {
	pushSender *push.PushSender
}

func NewPushController(pushSender *push.PushSender) *PushController {
	return &PushController{
		pushSender: pushSender,
	}
}

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
}

func (h *PushController) SendPush(c *gin.Context) {
	var req requests.SendPushNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Message: "Invalid request format",
		})
		return
	}

	notificationID := uuid.New().String()
	correlationID := c.GetHeader("X-Correlation-ID")
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	notification := &model.PushNotification{
		CorrelationID:      correlationID,
		NotificationID:     notificationID,
		PushToken:          req.PushToken,
		CreatedAt:          time.Now(),
		NotificattionTitle: req.Title,
		NotificattionBody:  req.Body,
	}

	err := h.pushSender.Send(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
			Message: "Failed to send push notification",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"notification_id": notificationID,
			"correlation_id":  correlationID,
			"status":          "sent",
			"timestamp":       time.Now().Unix(),
		},
		Message: "Push notification sent successfully",
	})
}
