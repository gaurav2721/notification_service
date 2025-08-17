package routes

import (
	"os"
	"strconv"

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

		// Setup user routes (controlled by feature flag)
		if isUserRoutesEnabled() {
			SetupUserRoutes(api, userHandler)
		}
	}
}

// isUserRoutesEnabled checks if user routes should be enabled based on environment variable
func isUserRoutesEnabled() bool {
	enableUserRoutes := os.Getenv("ENABLE_USER_ROUTES")
	if enableUserRoutes == "" {
		return false // Default to disabled
	}

	enabled, err := strconv.ParseBool(enableUserRoutes)
	if err != nil {
		return false // Default to disabled on parsing error
	}

	return enabled
}
