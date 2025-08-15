package notification_manager

import (
	"context"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager/templates"
)

// NotificationManagerImpl implements the NotificationManager interface
type NotificationManagerImpl struct {
	emailService    interface{}
	slackService    interface{}
	inappService    interface{}
	userService     interface{}
	scheduler       interface{}
	templateManager templates.TemplateManager
}

// NewNotificationManager creates a new notification manager instance
func NewNotificationManager(
	emailService interface{},
	slackService interface{},
	inappService interface{},
	userService interface{},
	scheduler interface{},
	templateManager templates.TemplateManager,
) *NotificationManagerImpl {
	return &NotificationManagerImpl{
		emailService:    emailService,
		slackService:    slackService,
		inappService:    inappService,
		userService:     userService,
		scheduler:       scheduler,
		templateManager: templateManager,
	}
}

// NewNotificationManagerWithDefaultTemplate creates a new notification manager with default template manager
func NewNotificationManagerWithDefaultTemplate(
	emailService interface{},
	slackService interface{},
	inappService interface{},
	userService interface{},
	scheduler interface{},
) *NotificationManagerImpl {
	return NewNotificationManager(
		emailService,
		slackService,
		inappService,
		userService,
		scheduler,
		templates.NewTemplateManager(),
	)
}

// SendNotification sends a notification through the appropriate channel
func (nm *NotificationManagerImpl) SendNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get notification
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    *models.TemplateData
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrUnsupportedNotificationType
	}

	// Check if it's a scheduled notification
	if notif.ScheduledAt != nil {
		return nm.ScheduleNotification(ctx, notification)
	}

	// Process template if provided
	if notif.Template != nil {
		// Validate template data
		if err := nm.templateManager.ValidateTemplateData(notif.Template.ID, notif.Template.Data); err != nil {
			return nil, err
		}

		// TODO: Process template variables and create content
		// For now, just return success
	}

	// Send notification based on type
	switch notif.Type {
	case "email":
		// TODO: Implement email sending
		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "sent",
			Message: "Email notification sent successfully",
			SentAt:  time.Now(),
			Channel: "email",
		}, nil
	case "slack":
		// TODO: Implement Slack sending
		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "sent",
			Message: "Slack notification sent successfully",
			SentAt:  time.Now(),
			Channel: "slack",
		}, nil
	case "in_app":
		// TODO: Implement in-app notification
		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "sent",
			Message: "In-app notification sent successfully",
			SentAt:  time.Now(),
			Channel: "in_app",
		}, nil
	default:
		return nil, ErrUnsupportedNotificationType
	}
}

// SendNotificationToUsers sends notifications to specific users
func (nm *NotificationManagerImpl) SendNotificationToUsers(ctx context.Context, userIDs []string, notification interface{}) (interface{}, error) {
	// Type assertion to get notification
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    *models.TemplateData
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrUnsupportedNotificationType
	}

	// Set recipients to the provided user IDs
	notif.Recipients = userIDs

	// Send notification
	return nm.SendNotification(ctx, notification)
}

// ScheduleNotification schedules a notification for future delivery
func (nm *NotificationManagerImpl) ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get notification
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    *models.TemplateData
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrNoScheduledTime
	}

	if notif.ScheduledAt == nil {
		return nil, ErrNoScheduledTime
	}

	// Schedule the notification
	// err := nm.scheduler.ScheduleJob(notif.ID, *notif.ScheduledAt, func() {
	// 	nm.SendNotification(context.Background(), notification)
	// })

	// if err != nil {
	// 	return nil, err
	// }

	// Return success response
	return &struct {
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		Message     string    `json:"message"`
		ScheduledAt time.Time `json:"scheduled_at"`
	}{
		ID:          notif.ID,
		Status:      "scheduled",
		Message:     "Notification scheduled successfully",
		ScheduledAt: *notif.ScheduledAt,
	}, nil
}

// GetNotificationStatus retrieves the status of a notification
func (nm *NotificationManagerImpl) GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error) {
	// This would typically query a database or storage system
	// For now, return a mock response
	return &struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     notificationID,
		Status: "sent",
	}, nil
}

// CreateTemplate creates a new notification template
func (nm *NotificationManagerImpl) CreateTemplate(ctx context.Context, template interface{}) (interface{}, error) {
	// Type assertion to get template
	tmpl, ok := template.(*models.Template)
	if !ok {
		return nil, ErrTemplateNotFound
	}

	return nm.templateManager.CreateTemplate(ctx, tmpl)
}

// GetTemplateVersion retrieves a specific version of a notification template
func (nm *NotificationManagerImpl) GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error) {
	return nm.templateManager.GetTemplateVersion(ctx, templateID, version)
}

// GetPredefinedTemplates returns all predefined templates
func (nm *NotificationManagerImpl) GetPredefinedTemplates() []*models.Template {
	return nm.templateManager.GetPredefinedTemplates()
}
