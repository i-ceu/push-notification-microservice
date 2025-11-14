package push

import (
	"fmt"
	"log"
	"push_service/config"
	"push_service/model"
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
	credPath, err := config.GetFirebaseCredentials()
	if err == nil && credPath != "" {
		log.Printf("Initializing FCM with credentials")

		serviceAccountSender, err := NewFCMSender(credPath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize FCM: %w", err)
		}

		ps.fcmSender = serviceAccountSender
		return ps, nil
	}

	return &PushSender{
		config:         cfg,
		fcmSender:      ps.fcmSender,
		circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
	}, nil
}

func (ps *PushSender) Send(notification *model.PushNotification) error {
	fmt.Println("Attempting send notification")
	return ps.circuitBreaker.Call(func() error {

		if ps.fcmSender == nil {
			return fmt.Errorf("FCM sender not configured")
		}
		return ps.fcmSender.Send(notification)

	})
}
