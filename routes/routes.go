package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/routes/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, notificationHandler *handlers.NotificationHandler, userHandler *handlers.UserHandler) {
	// Setup middleware
	middleware.SetupMiddleware(router)

	// Setup health routes
	SetupHealthRoutes(router, notificationHandler)

	// API routes
	api := router.Group("/api/v1")
	{
		// Setup notification routes
		SetupNotificationRoutes(api, notificationHandler)

		// Setup template routes
		SetupTemplateRoutes(api, notificationHandler)

		// Setup user routes
		SetupUserRoutes(api, userHandler)
	}
}
