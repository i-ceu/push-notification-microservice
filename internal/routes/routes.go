package routes

import (
	"log"
	"net/http"
	"push-notification-microservice/internal/config"
	"push-notification-microservice/internal/controllers"
	"push-notification-microservice/internal/queue"
	"push-notification-microservice/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initialize(cfg *config.Config) *http.Server {
	// logFile := helpers.SetupLogging()
	// defer logFile.Close()

	// gin.DefaultWriter = logFile

	pushSender, err := services.NewPushSender(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize push sender: %v", err)
	}

	consumer, err := queue.NewConsumer(cfg, pushSender)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	go consumer.Start()

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

	return &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

}
