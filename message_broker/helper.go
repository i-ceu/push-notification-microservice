package message_broker

import (
	"log"
)


func RabbitMqError (err error, message string) {
	if err != nil {
		log.Fatalf(`%s:%s`, message, err)
	}

}