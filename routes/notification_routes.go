package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/validation"
	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes configures notification-related routes
func SetupNotificationRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
	// Create validation layer
	validationLayer := validation.NewValidationLayer()

	// Notification endpoints with validation
	api.POST("/notifications", validationLayer.ValidateNotificationRequest(), handler.SendNotification)
	api.GET("/notifications/:id", validationLayer.ValidateNotificationID(), handler.GetNotificationStatus)
}
