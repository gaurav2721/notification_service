package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gaurav2721/notification-service/constants"
	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware validates API keys for protected routes
func APIKeyMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get API key from environment variable
		expectedAPIKey := os.Getenv(constants.API_KEY)

		if expectedAPIKey == "" {
			// If no API key is configured, allow the request (for development)
			c.Next()
			return
		}

		// Get API key from request header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "API key is required",
				"message": "Please provide an API key in the Authorization header",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer " or "ApiKey "
		var providedAPIKey string
		if strings.HasPrefix(authHeader, "Bearer ") {
			providedAPIKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "ApiKey ") {
			providedAPIKey = strings.TrimPrefix(authHeader, "ApiKey ")
		} else {
			// Try to use the header value directly as API key
			providedAPIKey = authHeader
		}

		// Validate API key
		if providedAPIKey != expectedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid API key",
				"message": "The provided API key is invalid",
			})
			c.Abort()
			return
		}

		// API key is valid, proceed
		c.Next()
	})
}
