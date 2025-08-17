package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gaurav2721/notification-service/constants"
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/logger"
	"github.com/gaurav2721/notification-service/routes"
	"github.com/gaurav2721/notification-service/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Error("No .env file found, using system environment variables")
		return
	}

	// Configure logging
	logger.Configure()

	// Initialize service container (manages all service dependencies)
	serviceContainer := services.NewServiceContainer()

	// Initialize handlers with required dependencies
	notificationHandler := handlers.NewNotificationHandler(serviceContainer.GetNotificationService(), serviceContainer.GetUserService(), serviceContainer.GetKafkaService())
	userHandler := handlers.NewUserHandler(serviceContainer.GetUserService())
	logrus.Debug("Handlers initialized successfully")

	// Setup Gin router
	router := gin.Default()

	// Setup all routes using the routes package
	routes.SetupRoutes(router, notificationHandler, userHandler)
	logrus.Debug("Routes configured successfully")

	// Get port from environment or use default
	port := os.Getenv(constants.PORT)
	if port == "" {
		port = constants.DefaultPort
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logrus.Info("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// Start server in a goroutine
	go func() {
		logrus.WithField("port", port).Debug("Starting notification service")
		if err := router.Run(":" + port); err != nil {
			logrus.WithError(err).Error("Server error")
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Gracefully shutdown services
	logrus.Debug("Initiating graceful shutdown of services")
	if err := serviceContainer.Shutdown(context.Background()); err != nil {
		logrus.WithError(err).Error("Error during shutdown")
	}

	logrus.Debug("Notification service stopped gracefully")
}
