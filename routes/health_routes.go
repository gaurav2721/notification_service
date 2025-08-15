package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(router *gin.Engine, handler *handlers.NotificationHandler) {
	// Health check endpoint
	router.GET("/health", handler.HealthCheck)
}
