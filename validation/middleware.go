package validation

import (
	"fmt"
	"net/http"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ValidationLayer provides middleware for request validation
type ValidationLayer struct {
	notificationValidator *NotificationValidator
	templateValidator     *TemplateValidator
}

// NewValidationLayer creates a new validation layer
func NewValidationLayer() *ValidationLayer {
	return &ValidationLayer{
		notificationValidator: NewNotificationValidator(),
		templateValidator:     NewTemplateValidator(),
	}
}

// ValidateNotificationRequest is middleware that validates notification requests
func (vm *ValidationLayer) ValidateNotificationRequest() gin.HandlerFunc {
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
func (vm *ValidationLayer) ValidateTemplateRequest() gin.HandlerFunc {
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
func (vm *ValidationLayer) ValidateTemplateID() gin.HandlerFunc {
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
func (vm *ValidationLayer) ValidateTemplateVersion() gin.HandlerFunc {
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

// ValidateNotificationID is middleware that validates notification ID parameter
func (vm *ValidationLayer) ValidateNotificationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		notificationID := c.Param("id")

		validationResult := vm.notificationValidator.ValidateNotificationID(notificationID)
		if !validationResult.IsValid {
			logrus.WithField("errors", validationResult.Errors).Warn("Validation failed for notification ID")
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
func (vm *ValidationLayer) ValidateUserRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic user validation can be added here
		// For now, just pass through
		c.Next()
	}
}
