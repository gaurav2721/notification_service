package notification_manager

import "context"

// NotificationManager interface defines methods for notification management
type NotificationManager interface {
	SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
	SendNotificationToUsers(ctx context.Context, userIDs []string, notification interface{}) (interface{}, error)
	ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
	GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
	CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
	GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
}
