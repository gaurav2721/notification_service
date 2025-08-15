package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures user-related routes
func SetupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	// User endpoints
	users := api.Group("/users")
	{
		users.GET("/", userHandler.GetUsers)         // Get all users
		users.GET("/:id", userHandler.GetUser)       // Get user by ID
		users.POST("/", userHandler.CreateUser)      // Create new user
		users.PUT("/:id", userHandler.UpdateUser)    // Update user
		users.DELETE("/:id", userHandler.DeleteUser) // Delete user

		// User notification specific endpoints
		users.GET("/:id/notification-info", userHandler.GetUserNotificationInfo) // Get user notification info

		// Device management endpoints
		users.POST("/:id/devices", userHandler.RegisterDevice)                        // Register a new device
		users.GET("/:id/devices", userHandler.GetUserDevices)                         // Get all devices for user
		users.GET("/:id/devices/active", userHandler.GetActiveUserDevices)            // Get active devices for user
		users.PUT("/devices/:deviceId", userHandler.UpdateDeviceInfo)                 // Update device information
		users.DELETE("/devices/:deviceId", userHandler.RemoveDevice)                  // Remove device
		users.PATCH("/devices/:deviceId/deactivate", userHandler.DeactivateDevice)    // Deactivate device
		users.PATCH("/devices/:deviceId/last-used", userHandler.UpdateDeviceLastUsed) // Update device last used
	}
}
