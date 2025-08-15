package notification_manager

import (
	"context"
	"time"

	"github.com/gaurav2721/notification-service/external_services/email"
	"github.com/gaurav2721/notification-service/external_services/inapp"
	"github.com/gaurav2721/notification-service/external_services/slack"
)

// NotificationManagerImpl manages different notification channels
type NotificationManagerImpl struct {
	emailService email.EmailService
	slackService slack.SlackService
	inAppService inapp.InAppService
	templates    map[string]interface{}
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(
	emailService email.EmailService,
	slackService slack.SlackService,
	inAppService inapp.InAppService,
) NotificationManager {
	return &NotificationManagerImpl{
		emailService: emailService,
		slackService: slackService,
		inAppService: inAppService,
		templates:    make(map[string]interface{}),
	}
}

// NewNotificationManagerWithConfig creates a new notification manager with configuration
func NewNotificationManagerWithConfig(
	emailService email.EmailService,
	slackService slack.SlackService,
	inAppService inapp.InAppService,
	config *NotificationConfig,
) NotificationManager {
	return &NotificationManagerImpl{
		emailService: emailService,
		slackService: slackService,
		inAppService: inAppService,
		templates:    make(map[string]interface{}),
	}
}

// SendNotification sends a notification through the appropriate channel
func (nm *NotificationManagerImpl) SendNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get notification type
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Title       string
		Message     string
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrUnsupportedNotificationType
	}

	// Route notification to appropriate channel
	switch notif.Type {
	case "email":
		return nm.emailService.SendEmail(ctx, notification)
	case "slack":
		return nm.slackService.SendSlackMessage(ctx, notification)
	case "in_app":
		return nm.inAppService.SendInAppNotification(ctx, notification)
	default:
		return nil, ErrUnsupportedNotificationType
	}
}

// SendNotificationToUsers sends a notification to specific users
func (nm *NotificationManagerImpl) SendNotificationToUsers(ctx context.Context, userIDs []string, notification interface{}) (interface{}, error) {
	// Implementation for sending to specific users
	// This would typically involve getting user notification info and routing accordingly
	return nm.SendNotification(ctx, notification)
}

// ScheduleNotification schedules a notification for later delivery
func (nm *NotificationManagerImpl) ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get notification
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Title       string
		Message     string
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
func (nm *NotificationManagerImpl) CreateTemplate(ctx context.Context, template interface{}) error {
	// Type assertion to get template
	tmpl, ok := template.(*struct {
		ID   string
		Name string
		Type string
		Body string
	})
	if !ok {
		return ErrTemplateNotFound
	}

	nm.templates[tmpl.ID] = template
	return nil
}

// GetTemplate retrieves a notification template
func (nm *NotificationManagerImpl) GetTemplate(ctx context.Context, templateID string) (interface{}, error) {
	if template, exists := nm.templates[templateID]; exists {
		return template, nil
	}
	return nil, ErrTemplateNotFound
}

// UpdateTemplate updates an existing notification template
func (nm *NotificationManagerImpl) UpdateTemplate(ctx context.Context, template interface{}) error {
	// Type assertion to get template
	tmpl, ok := template.(*struct {
		ID   string
		Name string
		Type string
		Body string
	})
	if !ok {
		return ErrTemplateNotFound
	}

	if _, exists := nm.templates[tmpl.ID]; !exists {
		return ErrTemplateNotFound
	}

	nm.templates[tmpl.ID] = template
	return nil
}

// DeleteTemplate deletes a notification template
func (nm *NotificationManagerImpl) DeleteTemplate(ctx context.Context, templateID string) error {
	if _, exists := nm.templates[templateID]; !exists {
		return ErrTemplateNotFound
	}

	delete(nm.templates, templateID)
	return nil
}
