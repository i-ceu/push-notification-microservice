package controllers

import (
	"net/http"
	"push-notification-microservice/internal/models"
	"push-notification-microservice/internal/requests"
	"push-notification-microservice/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PushController struct {
	pushSender *services.PushSender
}

func NewPushController(pushSender *services.PushSender) *PushController {
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

	notification := &models.PushNotification{
		CorrelationID:  correlationID,
		NotificationID: notificationID,
		PushToken:      req.PushToken,
		Title:          req.Title,
		Body:           req.Body,
		ImageURL:       req.ImageURL,
		ClickURL:       req.ClickURL,
		Data:           req.Data,
		Variables:      req.Variables,
		CreatedAt:      time.Now(),
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

func (h *PushController) SendBatchPush(c *gin.Context) {
	var req requests.BatchPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Message: "Invalid request format",
		})
		return
	}

	correlationID := c.GetHeader("X-Correlation-ID")
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	var successCount, failureCount int
	results := make([]map[string]any, 0)

	for _, token := range req.PushTokens {
		notificationID := uuid.New().String()

		notification := &models.PushNotification{
			CorrelationID:  correlationID,
			NotificationID: notificationID,
			PushToken:      token,
			Title:          req.Title,
			Body:           req.Body,
			ImageURL:       req.ImageURL,
			ClickURL:       req.ClickURL,
			Data:           req.Data,
			Variables:      req.Variables,
			CreatedAt:      time.Now(),
		}

		err := h.pushSender.Send(notification)

		result := map[string]any{
			"push_token":      token,
			"notification_id": notificationID,
		}

		if err != nil {
			failureCount++
			result["status"] = "failed"
			result["error"] = err.Error()
		} else {
			successCount++
			result["status"] = "sent"
		}

		results = append(results, result)
	}

	statusCode := http.StatusOK
	if successCount == 0 {
		statusCode = http.StatusInternalServerError
	} else if failureCount > 0 {
		statusCode = http.StatusMultiStatus
	}

	c.JSON(statusCode, Response{
		Success: successCount > 0,
		Data: gin.H{
			"correlation_id": correlationID,
			"total":          len(req.PushTokens),
			"success_count":  successCount,
			"failure_count":  failureCount,
			"results":        results,
			"timestamp":      time.Now().Unix(),
		},
		Message: "Batch push notification processing completed",
	})
}
