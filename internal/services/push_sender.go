package services

import (
	"fmt"
	"log"
	"push-notification-microservice/internal/config"
	"push-notification-microservice/internal/models"
	"time"
)

type PushSender struct {
	config         *config.Config
	fcmSender      *FCMSender
	circuitBreaker *CircuitBreaker
}

func NewPushSender(cfg *config.Config) (*PushSender, error) {
	ps := &PushSender{
		config:         cfg,
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
	}
	// Initialize FCM sender if configured
	credPath, err := config.GetFirebaseCredentials()
	if err == nil && credPath != "" {
		log.Printf("🔐 Initializing FCM with credentials")

		serviceAccountSender, err := NewFCMSender(credPath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize FCM: %w", err)
		}

		ps.fcmSender = serviceAccountSender
		log.Println("✅ FCM initialized successfully")
		return ps, nil
	}

	return &PushSender{
		config:         cfg,
		fcmSender:      ps.fcmSender,
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
	}, nil
}

func (ps *PushSender) Send(notification *models.PushNotification) error {
	return ps.circuitBreaker.Call(func() error {
		switch ps.config.PushProvider {
		case "fcm":
			if ps.fcmSender == nil {
				return fmt.Errorf("FCM sender not configured")
			}
			return ps.fcmSender.Send(notification)
		default:
			return fmt.Errorf("unknown push provider: %s", ps.config.PushProvider)
		}
	})
}
