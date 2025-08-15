package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes configures notification-related routes
func SetupNotificationRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
	// Notification endpoints
	api.POST("/notifications", handler.SendNotification)
	api.GET("/notifications/:id", handler.GetNotificationStatus)
}
