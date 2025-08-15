package notification

import "context"

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
