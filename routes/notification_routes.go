package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/middleware"
	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes configures notification-related routes
func SetupNotificationRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
	// Create validation middleware
	validationMiddleware := middleware.NewValidationMiddleware()

	// Notification endpoints with validation
	api.POST("/notifications", validationMiddleware.ValidateNotificationRequest(), handler.SendNotification)
	api.GET("/notifications/:id", handler.GetNotificationStatus)
}
