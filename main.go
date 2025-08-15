package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/routes"
	"github.com/gaurav2721/notification-service/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize service container (manages all service dependencies)
	serviceContainer := services.NewServiceContainer()

	// Initialize handlers with interface dependencies from the container
	notificationHandler := handlers.NewNotificationHandler(serviceContainer.GetNotificationService())
	userHandler := handlers.NewUserHandler(serviceContainer.GetUserService())

	// Setup Gin router
	router := gin.Default()

	// Setup all routes using the routes package
	routes.SetupRoutes(router, notificationHandler, userHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// Start server in a goroutine
	go func() {
		log.Printf("Starting notification service on port %s", port)
		if err := router.Run(":" + port); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Gracefully shutdown services
	if err := serviceContainer.Shutdown(context.Background()); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Notification service stopped gracefully")
}
