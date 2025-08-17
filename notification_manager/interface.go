package notification_manager

import (
	"context"

	"github.com/gaurav2721/notification-service/models"
)

// NotificationManager interface defines methods for notification management
type NotificationManager interface {
	SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
	SendNotificationToUsers(ctx context.Context, userIDs []string, notification interface{}) (interface{}, error)
	ScheduleNotification(ctx context.Context, notificationId string, notification *models.NotificationRequest, job func() error) error

	GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
	SetNotificationStatus(ctx context.Context, notificationId string, notification *models.NotificationRequest, status string) error
	CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
	GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
	GetTemplateByID(templateID string) (*models.Template, error)
	GetTemplateByIDAndVersion(templateID string, version int) (*models.Template, error)
	GetPredefinedTemplates() []*models.Template
}
