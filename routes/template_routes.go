package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/validation"
	"github.com/gin-gonic/gin"
)

// SetupTemplateRoutes configures template-related routes
func SetupTemplateRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
	// Create validation layer
	validationLayer := validation.NewValidationLayer()

	// Template endpoints with validation
	api.POST("/templates", validationLayer.ValidateTemplateRequest(), handler.CreateTemplate)
	api.GET("/templates/predefined", handler.GetPredefinedTemplates)
	api.GET("/templates/:templateId/versions/:version",
		validationLayer.ValidateTemplateID(),
		validationLayer.ValidateTemplateVersion(),
		handler.GetTemplateVersion)
}
