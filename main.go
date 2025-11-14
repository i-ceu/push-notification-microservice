package main

import (
	"log"
	"net/http"
	"push_service/config"
	"push_service/message_broker"
	"push_service/push"
	"push_service/routes"
)

func main() {

	cfg := config.Load()

	pushSender, err := push.NewPushSender(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize push sender: %v", err)
	}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Push service is running"))
		})
		routes := routes.RegisterRoutes(cfg, pushSender)

		log.Printf("Push running on port %s", cfg.Port)
		if err := http.ListenAndServe(cfg.Port, routes); err != nil {
			log.Fatal(err)
		}

	}()

	// Connect to RabbitMQ
	connection, channel := message_broker.Connect_to_rabitmq()

	defer connection.Close()
	defer channel.Close()

	// Consumer
	message_broker.Consumer(channel, pushSender)

}
