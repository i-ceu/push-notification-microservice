package message_broker

import (
	"log"
	"os"


	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect_to_rabitmq () (*amqp.Connection, *amqp.Channel){
    log.Println("Trying to connect to Rabbitmq")

	Url := os.Getenv("RABBITMQ_URL")
	


	connection, err := amqp.Dial(Url)
	 RabbitMqError(err, "unable to connect to rabbitmq")
	 
	log.Println("Rabbitmq connected")


		  
	 channel , err := connection.Channel()
	  RabbitMqError(err, "failed to connect to channel")

	  log.Println("channel registed")




    
return connection, channel

}