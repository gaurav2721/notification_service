package middleware

import (
	"net/http"

	"github.com/gaurav2721/notification-service/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ValidationMiddleware provides middleware for request validation
type ValidationMiddleware struct {
	validator *validation.NotificationValidator
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validation.NewNotificationValidator(),
	}
}

// ValidateNotificationRequest is middleware that validates notification requests
func (vm *ValidationMiddleware) ValidateNotificationRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request validation.NotificationRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			logrus.WithError(err).Warn("Invalid JSON in notification request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		validationResult := vm.validator.ValidateNotificationRequest(&request)
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

// ValidateTemplateRequest is middleware that validates template requests
func (vm *ValidationMiddleware) ValidateTemplateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic template validation can be added here
		// For now, just pass through
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
