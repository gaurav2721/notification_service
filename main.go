package main

import (
	"log"
	"os"

	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize services
	emailService := services.NewEmailService()
	slackService := services.NewSlackService()
	inAppService := services.NewInAppService()
	schedulerService := services.NewSchedulerService()

	// Initialize notification manager
	notificationManager := services.NewNotificationManager(
		emailService,
		slackService,
		inAppService,
		schedulerService,
	)

	// Initialize handler
	handler := handlers.NewNotificationHandler(notificationManager)

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Notification endpoints
		api.POST("/notifications", handler.SendNotification)
		api.GET("/notifications/:id", handler.GetNotificationStatus)

		// Template endpoints
		api.POST("/templates", handler.CreateTemplate)
		api.GET("/templates/:id", handler.GetTemplate)
		api.PUT("/templates/:id", handler.UpdateTemplate)
		api.DELETE("/templates/:id", handler.DeleteTemplate)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting notification service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
