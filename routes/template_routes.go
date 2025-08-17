package routes

import (
	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/validation"
	"github.com/gin-gonic/gin"
)

// SetupTemplateRoutes configures template-related routes
func SetupTemplateRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
	// Create validation middleware
	validationMiddleware := validation.NewValidationMiddleware()

	// Template endpoints with validation
	api.POST("/templates", validationMiddleware.ValidateTemplateRequest(), handler.CreateTemplate)
	api.GET("/templates/predefined", handler.GetPredefinedTemplates)
	api.GET("/templates/:templateId/versions/:version",
		validationMiddleware.ValidateTemplateID(),
		validationMiddleware.ValidateTemplateVersion(),
		handler.GetTemplateVersion)
}
