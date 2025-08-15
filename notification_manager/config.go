package notification_manager

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
