package notification

import (
	"context"
	"errors"
)

// NotificationService interface defines methods for notification management
type NotificationService interface {
	SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
	SendNotificationToUsers(ctx context.Context, userIDs []string, notification interface{}) (interface{}, error)
	ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
	GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
	CreateTemplate(ctx context.Context, template interface{}) error
	GetTemplate(ctx context.Context, templateID string) (interface{}, error)
	UpdateTemplate(ctx context.Context, template interface{}) error
	DeleteTemplate(ctx context.Context, templateID string) error
}

// NotificationConfig holds notification service configuration
type NotificationConfig struct {
	DefaultPriority string
	MaxRetries      int
	RetryDelay      int
}

// DefaultNotificationConfig returns default notification configuration
func DefaultNotificationConfig() *NotificationConfig {
	return &NotificationConfig{
		DefaultPriority: "normal",
		MaxRetries:      3,
		RetryDelay:      1000,
	}
}

// Notification service errors
var (
	ErrUnsupportedNotificationType = errors.New("unsupported notification type")
	ErrNoScheduledTime             = errors.New("no scheduled time provided")
	ErrTemplateNotFound            = errors.New("template not found")
	ErrInvalidRecipients           = errors.New("invalid recipients")
)
