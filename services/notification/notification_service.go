package notification

import (
	"context"
	"time"

	"github.com/gaurav2721/notification-service/services/common"
)

// NotificationService defines the interface for notification operations
type NotificationService interface {
	SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
	ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
	GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
	CreateTemplate(ctx context.Context, template interface{}) error
	GetTemplate(ctx context.Context, templateID string) (interface{}, error)
	UpdateTemplate(ctx context.Context, template interface{}) error
	DeleteTemplate(ctx context.Context, templateID string) error
}

// NotificationManager manages different notification channels
type NotificationManager struct {
	emailService common.EmailService
	slackService common.SlackService
	inAppService common.InAppService
	scheduler    common.SchedulerService
	templates    map[string]interface{}
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(emailService common.EmailService, slackService common.SlackService, inAppService common.InAppService, scheduler common.SchedulerService) *NotificationManager {
	return &NotificationManager{
		emailService: emailService,
		slackService: slackService,
		inAppService: inAppService,
		scheduler:    scheduler,
		templates:    make(map[string]interface{}),
	}
}

// SendNotification sends a notification through the appropriate channel
func (nm *NotificationManager) SendNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get notification type
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Priority    string
		Title       string
		Message     string
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, common.ErrUnsupportedNotificationType
	}

	// Check if notification is scheduled for future
	if notif.ScheduledAt != nil && notif.ScheduledAt.After(time.Now()) {
		return nm.ScheduleNotification(ctx, notification)
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
		return nil, common.ErrUnsupportedNotificationType
	}
}

// ScheduleNotification schedules a notification for later delivery
func (nm *NotificationManager) ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get notification
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Priority    string
		Title       string
		Message     string
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, common.ErrNoScheduledTime
	}

	if notif.ScheduledAt == nil {
		return nil, common.ErrNoScheduledTime
	}

	// Schedule the notification
	err := nm.scheduler.ScheduleJob(notif.ID, *notif.ScheduledAt, func() {
		nm.SendNotification(context.Background(), notification)
	})

	if err != nil {
		return nil, common.ErrSchedulingFailed
	}

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
func (nm *NotificationManager) GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error) {
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
func (nm *NotificationManager) CreateTemplate(ctx context.Context, template interface{}) error {
	// Type assertion to get template
	tmpl, ok := template.(*struct {
		ID   string
		Name string
		Type string
		Body string
	})
	if !ok {
		return common.ErrTemplateNotFound
	}

	nm.templates[tmpl.ID] = template
	return nil
}

// GetTemplate retrieves a notification template
func (nm *NotificationManager) GetTemplate(ctx context.Context, templateID string) (interface{}, error) {
	if template, exists := nm.templates[templateID]; exists {
		return template, nil
	}
	return nil, common.ErrTemplateNotFound
}

// UpdateTemplate updates an existing notification template
func (nm *NotificationManager) UpdateTemplate(ctx context.Context, template interface{}) error {
	// Type assertion to get template
	tmpl, ok := template.(*struct {
		ID   string
		Name string
		Type string
		Body string
	})
	if !ok {
		return common.ErrTemplateNotFound
	}

	if _, exists := nm.templates[tmpl.ID]; !exists {
		return common.ErrTemplateNotFound
	}

	nm.templates[tmpl.ID] = template
	return nil
}

// DeleteTemplate deletes a notification template
func (nm *NotificationManager) DeleteTemplate(ctx context.Context, templateID string) error {
	if _, exists := nm.templates[templateID]; !exists {
		return common.ErrTemplateNotFound
	}

	delete(nm.templates, templateID)
	return nil
}
