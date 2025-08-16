package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/routes"
	"github.com/gaurav2721/notification-service/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configure logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Set log level from environment variable or default to InfoLevel
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel) // Default to InfoLevel to disable debug logs
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using system environment variables")
	}

	logrus.Info("Starting notification service initialization")

	// Initialize service container (manages all service dependencies)
	serviceContainer := services.NewServiceContainer()
	logrus.Info("Service container initialized successfully")

	// Initialize handlers with required dependencies
	notificationHandler := handlers.NewNotificationHandler(serviceContainer.GetNotificationService(), serviceContainer.GetUserService(), serviceContainer.GetKafkaService())
	userHandler := handlers.NewUserHandler(serviceContainer.GetUserService())
	logrus.Info("Handlers initialized successfully")

	// Setup Gin router
	router := gin.Default()

	// Setup all routes using the routes package
	routes.SetupRoutes(router, notificationHandler, userHandler)
	logrus.Info("Routes configured successfully")

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
		logrus.Info("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// Start server in a goroutine
	go func() {
		logrus.WithField("port", port).Info("Starting notification service")
		if err := router.Run(":" + port); err != nil {
			logrus.WithError(err).Error("Server error")
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Gracefully shutdown services
	logrus.Info("Initiating graceful shutdown of services")
	if err := serviceContainer.Shutdown(context.Background()); err != nil {
		logrus.WithError(err).Error("Error during shutdown")
	}

	logrus.Info("Notification service stopped gracefully")
}
