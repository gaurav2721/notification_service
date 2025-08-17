package middleware

import (
	"github.com/gin-gonic/gin"
)

// SetupMiddleware configures all middleware for the application
func SetupMiddleware(router *gin.Engine) {
	// Add logging middleware
	router.Use(gin.Logger())

	// Add recovery middleware
	router.Use(gin.Recovery())
}
