package middleware

import (
	"fmt"
	"net/http"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ValidationMiddleware provides middleware for request validation
type ValidationMiddleware struct {
	notificationValidator *validation.NotificationValidator
	templateValidator     *validation.TemplateValidator
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		notificationValidator: validation.NewNotificationValidator(),
		templateValidator:     validation.NewTemplateValidator(),
	}
}

// ValidateNotificationRequest is middleware that validates notification requests
func (vm *ValidationMiddleware) ValidateNotificationRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.NotificationRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			logrus.WithError(err).Warn("Invalid JSON in notification request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		validationResult := vm.notificationValidator.ValidateNotificationRequest(&request)
		if !validationResult.IsValid {
			logrus.WithField("errors", validationResult.Errors).Warn("Validation failed for notification request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult.Errors,
			})
			c.Abort()
			return
		}

		// Store validated request in context for later use
		c.Set("validated_request", &request)
		c.Next()
	}
}

// ValidateTemplateRequest is middleware that validates template creation requests
func (vm *ValidationMiddleware) ValidateTemplateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.TemplateRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			logrus.WithError(err).Warn("Invalid JSON in template request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		validationResult := vm.templateValidator.ValidateTemplateRequest(&request)
		if !validationResult.IsValid {
			logrus.WithField("errors", validationResult.Errors).Warn("Validation failed for template request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult.Errors,
			})
			c.Abort()
			return
		}

		// Store validated request in context for later use
		c.Set("validated_template_request", &request)
		c.Next()
	}
}

// ValidateTemplateID is middleware that validates template ID parameter
func (vm *ValidationMiddleware) ValidateTemplateID() gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("templateId")

		validationResult := vm.templateValidator.ValidateTemplateID(templateID)
		if !validationResult.IsValid {
			logrus.WithField("errors", validationResult.Errors).Warn("Validation failed for template ID")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult.Errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateTemplateVersion is middleware that validates template version parameter
func (vm *ValidationMiddleware) ValidateTemplateVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		versionStr := c.Param("version")
		if versionStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "version parameter is required",
			})
			c.Abort()
			return
		}

		// Convert string to int
		version := 0
		if _, err := fmt.Sscanf(versionStr, "%d", &version); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "version must be a valid integer",
			})
			c.Abort()
			return
		}

		validationResult := vm.templateValidator.ValidateTemplateVersion(version)
		if !validationResult.IsValid {
			logrus.WithField("errors", validationResult.Errors).Warn("Validation failed for template version")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationResult.Errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateUserRequest is middleware that validates user requests
func (vm *ValidationMiddleware) ValidateUserRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic user validation can be added here
		// For now, just pass through
		c.Next()
	}
}
