package notification_manager

import (
	"context"

	"github.com/gaurav2721/notification-service/models"
)

// NotificationManager interface defines methods for notification management
type NotificationManager interface {
	GetNotificationStatus(notificationID string) (interface{}, error)
	CreateTemplate(template interface{}) (interface{}, error)
	GetTemplateVersion(templateID string, version int) (interface{}, error)
	GetPredefinedTemplates() []*models.Template

	// Main method for handling complete notification processing
	ProcessNotificationRequest(ctx context.Context, request *models.NotificationRequest) (interface{}, error)
}
