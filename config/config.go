package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL         string
	QueueName           string
	FCM_SERVICE_ACCOUNT string
	PushProvider        string
	Port                string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	credPath, _ := GetFirebaseCredentials()

	return &Config{
		RabbitMQURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:           getEnv("QUEUE_NAME", "push.queue"),
		FCM_SERVICE_ACCOUNT: getEnv("FCM_SERVICE_ACCOUNT", credPath),
		PushProvider:        getEnv("PUSH_PROVIDER", "fcm"),
		Port:                getEnv("SERVER_PORT", ":8084"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
