package inapp

import (
	"context"
	"errors"
)

// InAppService interface defines methods for in-app notifications
type InAppService interface {
	SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error)
}

// InAppConfig holds in-app service configuration
type InAppConfig struct {
	DatabaseURL string
	MaxRetries  int
	RetryDelay  int
}

// DefaultInAppConfig returns default in-app configuration
func DefaultInAppConfig() *InAppConfig {
	return &InAppConfig{
		DatabaseURL: "memory://",
		MaxRetries:  3,
		RetryDelay:  1000,
	}
}

// InApp service errors
var (
	ErrInAppSendFailed     = errors.New("failed to send in-app notification")
	ErrInAppDeviceToken    = errors.New("invalid in-app device token")
	ErrInAppDeviceNotFound = errors.New("in-app device not found")
)
