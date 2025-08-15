package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gin-gonic/gin"
)

// SetupTemplateRoutes configures template-related routes
func SetupTemplateRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
	// Template endpoints
	api.POST("/templates", handler.CreateTemplate)
	api.GET("/templates/:templateId/versions/:version", handler.GetTemplateVersion)
}
