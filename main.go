package main

import (
	"log"
	"os"

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

	// Initialize services
	emailService := services.NewEmailService()
	slackService := services.NewSlackService()
	inAppService := services.NewInAppService()
	schedulerService := services.NewSchedulerService()
	userService := services.NewUserService()

	// Initialize notification manager
	notificationManager := services.NewNotificationManager(
		emailService,
		slackService,
		inAppService,
		schedulerService,
	)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(notificationManager)
	userHandler := handlers.NewUserHandler(userService)

	// Setup Gin router
	router := gin.Default()

	// Setup all routes using the routes package
	routes.SetupRoutes(router, notificationHandler, userHandler)

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
