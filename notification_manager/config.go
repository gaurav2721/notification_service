package notification_manager

// NotificationConfig holds notification service configuration
type NotificationConfig struct {
	MaxRetries int
	RetryDelay int
}

// DefaultNotificationConfig returns default notification configuration
func DefaultNotificationConfig() *NotificationConfig {
	return &NotificationConfig{
		MaxRetries: 3,
		RetryDelay: 1000,
	}
}
