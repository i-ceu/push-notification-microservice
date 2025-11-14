package routes

import (
	"push_service/config"
	"push_service/controllers"
	"push_service/push"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(cfg *config.Config, pushSender *push.PushSender) *gin.Engine {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	healthHandler := controllers.NewHealthController()
	notificationHandler := controllers.NewPushController(pushSender)

	router.GET("/health", healthHandler.Check)

	api := router.Group("/api/v1")
	{
		push := api.Group("/push")
		{
			push.POST("/send", notificationHandler.SendPush)
		}
	}

	return router
}
