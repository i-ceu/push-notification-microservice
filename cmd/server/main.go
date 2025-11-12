package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"push-notification-microservice/internal/config"
	"push-notification-microservice/internal/routes"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	srv := routes.Initialize(cfg)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
