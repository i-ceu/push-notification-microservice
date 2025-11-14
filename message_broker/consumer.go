package message_broker

import (
	"encoding/json"
	"fmt"
	"log"
	"push_service/model"
	"push_service/push"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consumer(channel *amqp.Channel, pushSender *push.PushSender) {
	messages, err := channel.Consume(
		"push.queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println(err)
	}

	log.Println("notification consumed")

	forever := make(chan bool)

	go func() {
		for data := range messages {

			var messagequed model.PushNotification

			err := json.Unmarshal([]byte(data.Body), &messagequed)
			if err != nil {
				log.Println("Error decoding json", err)
				continue
			}
			err = pushSender.Send(&messagequed)
			fmt.Println("error sending message:", err)
			if err == nil {
				data.Ack(false)
				log.Println("message sent succesfully")
				break
			} else {
				channel.Publish(
					"",
					"failed.queue",
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        data.Body,
					},
				)
				log.Println("Failed to send message, published to failed queue")
				//publish to failed queue
			}

		}

	}()

	<-forever

}
