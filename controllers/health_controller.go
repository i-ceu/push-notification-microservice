package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	startTime time.Time
}

func NewHealthController() *HealthController {
	return &HealthController{
		startTime: time.Now(),
	}
}

func (h *HealthController) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "push-service",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(h.startTime).String(),
	})
}
