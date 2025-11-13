package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"push-notification-microservice/internal/config"

	"push-notification-microservice/internal/models"

	"push-notification-microservice/internal/services"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	pushSender *services.PushSender
	maxRetries int
}

func NewConsumer(cfg *config.Config, pushSender *services.PushSender) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel failed: %w", err)
	}

	if err := ch.Qos(1, 0, false); err != nil {
		return nil, fmt.Errorf("qos failed: %w", err)
	}

	q, err := ch.QueueDeclarePassive(
		cfg.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare failed: %w", err)
	}

	return &Consumer{
		conn:       conn,
		channel:    ch,
		queue:      q,
		pushSender: pushSender,
		maxRetries: 3,
	}, nil
}

func (c *Consumer) Start() {
	consumertag := fmt.Sprintf("consumer-%d", time.Now().Unix())

	msgs, err := c.channel.Consume(
		c.queue.Name,
		consumertag, // consumer tag
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	log.Printf("Consumer started: %s", consumertag)

	fmt.Println(len(msgs))
	for msg := range msgs {
		c.processMessage(msg)
	}
}

func (c *Consumer) processMessage(msg amqp.Delivery) {

	var notification models.PushNotification
	err := json.Unmarshal(msg.Body, &notification)
	if err != nil {
		log.Printf("[%s] Failed to unmarshal message: %v", notification.CorrelationID, err)
		msg.Ack(false)
		return
	}

	log.Printf("[%s] Processing push notification for token %s (attempt %d)",
		notification.CorrelationID,
		notification.PushToken[:min(10, len(notification.PushToken))]+"...",
		notification.RetryCount+1)

	if notification.PushToken == "" {
		log.Printf("[%s] Invalid push token - skipping", notification.CorrelationID)
		msg.Ack(false)
		return
	}

	err = c.pushSender.Send(&notification)
	if err != nil {
		log.Printf("Push notification sending failed: %v", err)
		c.handleFailure(msg, notification)
		return
	}

	log.Printf("Push notification sent successfully")
	msg.Ack(false)
}

func (c *Consumer) handleFailure(msg amqp.Delivery, notification models.PushNotification) {
	notification.RetryCount++

	if notification.RetryCount >= c.maxRetries {
		log.Printf("[%s] Max retries reached. Moving to dead-letter queue", notification.CorrelationID)
		c.sendToDeadLetterQueue(notification)
		msg.Ack(false)
		return
	}

	delay := time.Duration(1<<uint(notification.RetryCount)) * time.Second
	log.Printf("[%s] Retrying after %v", notification.CorrelationID, delay)

	time.Sleep(delay)
	msg.Nack(false, true)
}

func (c *Consumer) sendToDeadLetterQueue(notification models.PushNotification) {
	body, _ := json.Marshal(notification)
	err := c.channel.Publish(
		"",             // exchange
		"failed.queue", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("[%s] Failed to send to dead-letter queue: %v", notification.CorrelationID, err)
	}
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
