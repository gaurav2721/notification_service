package notification_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/gaurav2721/notification-service/external_services/kafka"
	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager/scheduler"
	"github.com/gaurav2721/notification-service/notification_manager/templates"
	"github.com/sirupsen/logrus"
)

// NotificationManagerImpl implements the NotificationManager interface
type NotificationManagerImpl struct {
	userService     interface{}
	kafkaService    kafka.KafkaService
	scheduler       scheduler.Scheduler
	templateManager templates.TemplateManager
	storage         *InMemoryStorage
}

// NewNotificationManager creates a new notification manager instance
func NewNotificationManager(
	userService interface{},
	kafkaService kafka.KafkaService,
	scheduler scheduler.Scheduler,
	templateManager templates.TemplateManager,
) *NotificationManagerImpl {
	return &NotificationManagerImpl{
		userService:     userService,
		kafkaService:    kafkaService,
		scheduler:       scheduler,
		templateManager: templateManager,
		storage:         NewInMemoryStorage(),
	}
}

// NewNotificationManagerWithDefaultTemplate creates a new notification manager with default template manager
// The scheduler is initialized internally within the notification manager
func NewNotificationManagerWithDefaultTemplate(
	userService interface{},
	kafkaService kafka.KafkaService,
) *NotificationManagerImpl {
	return NewNotificationManager(
		userService,
		kafkaService,
		scheduler.NewScheduler(),
		templates.NewTemplateManager(),
	)
}

// SendNotification sends a notification through the appropriate channel
func (nm *NotificationManagerImpl) SendNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	logrus.Debug("Starting notification send process")

	// Type assertion to get notification
	notif, ok := notification.(*struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    *models.TemplateData
		Recipients  []string
		ScheduledAt *time.Time
		From        *struct {
			Email string `json:"email"`
		}
	})
	if !ok {
		logrus.Error("Unsupported notification type")
		return nil, ErrUnsupportedNotificationType
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": notif.ID,
		"type":            notif.Type,
		"recipients":      len(notif.Recipients),
	}).Debug("Processing notification")

	// Check if it's a scheduled notification
	if notif.ScheduledAt != nil {
		logrus.Debug("Notification is scheduled, redirecting to scheduler")
		// Convert to NotificationRequest and schedule
		notificationRequest := &models.NotificationRequest{
			Type:        notif.Type,
			Content:     notif.Content,
			Template:    notif.Template,
			Recipients:  notif.Recipients,
			ScheduledAt: notif.ScheduledAt,
			From:        notif.From,
		}

		err := nm.ScheduleNotification(ctx, notif.ID, notificationRequest, func() error {
			// This job will be executed at the scheduled time
			logrus.WithField("notification_id", notif.ID).Info("Executing scheduled notification job")
			// TODO: Implement actual notification sending logic
			return nil
		})

		if err != nil {
			return nil, err
		}

		// Store the scheduled notification in memory
		if err := nm.storage.StoreNotification(notif.ID, notificationRequest); err != nil {
			logrus.WithError(err).WithField("notification_id", notif.ID).Warn("Failed to store scheduled notification")
		} else {
			// Set status to scheduled
			if err := nm.storage.UpdateNotificationStatus(notif.ID, StatusScheduled, ""); err != nil {
				logrus.WithError(err).WithField("notification_id", notif.ID).Warn("Failed to update scheduled notification status")
			}
		}

		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "scheduled",
			Message: "Notification scheduled successfully",
			SentAt:  time.Now(),
			Channel: notif.Type,
		}, nil
	}

	// Process template if provided
	if notif.Template != nil {
		logrus.Debug("Processing notification template")
		// Validate template data
		if err := nm.templateManager.ValidateTemplateData(notif.Template.ID, notif.Template.Data); err != nil {
			logrus.WithError(err).Error("Template validation failed")
			return nil, err
		}
		logrus.Debug("Template validation successful")

		// TODO: Process template variables and create content
		// For now, just return success
	}

	// Send notification to Kafka channel based on type
	logrus.WithField("type", notif.Type).Debug("Sending notification to Kafka channel")

	switch notif.Type {
	case "email":
		// Send to email channel
		if nm.kafkaService != nil {
			// Convert to JSON string (simplified for now)
			// TODO: Use proper JSON marshaling
			message := fmt.Sprintf("Email notification: %s", notif.ID)

			// Send to Kafka channel
			select {
			case nm.kafkaService.GetEmailChannel() <- message:
				logrus.WithField("notification_id", notif.ID).Debug("Email notification sent to Kafka channel")
				// Message sent successfully
			default:
				logrus.WithField("notification_id", notif.ID).Warn("Email channel is full")
				// Channel is full, handle accordingly
				// TODO: Add proper error handling
			}
		}

		// Store the email notification in memory
		emailNotificationRequest := &models.NotificationRequest{
			Type:       notif.Type,
			Content:    notif.Content,
			Template:   notif.Template,
			Recipients: notif.Recipients,
			From:       notif.From,
		}

		if err := nm.storage.StoreNotification(notif.ID, emailNotificationRequest); err != nil {
			logrus.WithError(err).WithField("notification_id", notif.ID).Warn("Failed to store email notification")
		} else {
			// Set status to queued
			if err := nm.storage.UpdateNotificationStatus(notif.ID, StatusQueued, ""); err != nil {
				logrus.WithError(err).WithField("notification_id", notif.ID).Warn("Failed to update email notification status")
			}
		}

		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "queued",
			Message: "Email notification queued for processing",
			SentAt:  time.Now(),
			Channel: "email",
		}, nil

	case "slack":
		// Send to slack channel
		if nm.kafkaService != nil {
			// Convert to JSON string (simplified for now)
			// TODO: Use proper JSON marshaling
			message := fmt.Sprintf("Slack notification: %s", notif.ID)

			// Send to Kafka channel
			select {
			case nm.kafkaService.GetSlackChannel() <- message:
				// Message sent successfully
			default:
				// Channel is full, handle accordingly
				// TODO: Add proper error handling
			}
		}

		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "queued",
			Message: "Slack notification queued for processing",
			SentAt:  time.Now(),
			Channel: "slack",
		}, nil

	case "ios_push":
		// Send to iOS push notification channel
		if nm.kafkaService != nil {
			// Convert to JSON string (simplified for now)
			// TODO: Use proper JSON marshaling
			message := fmt.Sprintf("iOS push notification: %s", notif.ID)

			// Send to Kafka channel
			select {
			case nm.kafkaService.GetIOSPushNotificationChannel() <- message:
				// Message sent successfully
			default:
				// Channel is full, handle accordingly
				// TODO: Add proper error handling
			}
		}

		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "queued",
			Message: "iOS push notification queued for processing",
			SentAt:  time.Now(),
			Channel: "ios_push",
		}, nil

	case "android_push":
		// Send to Android push notification channel
		if nm.kafkaService != nil {
			// Convert to JSON string (simplified for now)
			// TODO: Use proper JSON marshaling
			message := fmt.Sprintf("Android push notification: %s", notif.ID)

			// Send to Kafka channel
			select {
			case nm.kafkaService.GetAndroidPushNotificationChannel() <- message:
				// Message sent successfully
			default:
				// Channel is full, handle accordingly
				// TODO: Add proper error handling
			}
		}

		return &models.NotificationResponse{
			ID:      notif.ID,
			Status:  "queued",
			Message: "Android push notification queued for processing",
			SentAt:  time.Now(),
			Channel: "android_push",
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
func (nm *NotificationManagerImpl) ScheduleNotification(ctx context.Context, notificationId string, notification *models.NotificationRequest, job func() error) error {
	if notification == nil {
		return ErrUnsupportedNotificationType
	}

	if notification.ScheduledAt == nil {
		return ErrNoScheduledTime
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": notificationId,
		"scheduled_at":    notification.ScheduledAt,
		"type":            notification.Type,
		"recipients":      len(notification.Recipients),
	}).Debug("Scheduling notification job")

	// Convert the job function to match scheduler interface (func() instead of func() error)
	schedulerJob := func() {
		logrus.WithField("notification_id", notificationId).Info("Executing scheduled notification job")
		if err := job(); err != nil {
			logrus.WithError(err).WithField("notification_id", notificationId).Error("Scheduled notification job failed")
		} else {
			logrus.WithField("notification_id", notificationId).Info("Scheduled notification job completed successfully")
		}
	}

	// Schedule the job using the scheduler
	err := nm.scheduler.ScheduleJob(notificationId, *notification.ScheduledAt, schedulerJob)
	if err != nil {
		logrus.WithError(err).WithField("notification_id", notificationId).Error("Failed to schedule notification job")
		return err
	}

	logrus.WithField("notification_id", notificationId).Info("Notification job scheduled successfully")
	return nil
}

// GetNotificationStatus retrieves the status of a notification
func (nm *NotificationManagerImpl) GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error) {
	// Get notification from in-memory storage
	record, err := nm.storage.GetNotification(notificationID)
	if err != nil {
		logrus.WithError(err).WithField("notification_id", notificationID).Debug("Notification not found in storage")
		return nil, err
	}

	return &struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     record.ID,
		Status: string(record.Status),
	}, nil
}

// SetNotificationStatus sets the status of a notification
func (nm *NotificationManagerImpl) SetNotificationStatus(ctx context.Context, notificationId string, notification *models.NotificationRequest, status string) error {
	if notification == nil {
		return ErrUnsupportedNotificationType
	}

	if notificationId == "" {
		return fmt.Errorf("notification ID cannot be empty")
	}

	if status == "" {
		return fmt.Errorf("status cannot be empty")
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": notificationId,
		"status":          status,
		"type":            notification.Type,
		"recipients":      len(notification.Recipients),
	}).Debug("Setting notification status")

	// Convert string status to NotificationStatus type
	var notificationStatus NotificationStatus
	switch status {
	case "pending":
		notificationStatus = StatusPending
	case "scheduled":
		notificationStatus = StatusScheduled
	case "queued":
		notificationStatus = StatusQueued
	case "sent":
		notificationStatus = StatusSent
	case "failed":
		notificationStatus = StatusFailed
	case "cancelled":
		notificationStatus = StatusCancelled
	default:
		return fmt.Errorf("invalid status: %s", status)
	}

	// Store notification if it doesn't exist
	existingRecord, err := nm.storage.GetNotification(notificationId)
	if err != nil {
		// Notification doesn't exist, store it first
		if err := nm.storage.StoreNotification(notificationId, notification); err != nil {
			logrus.WithError(err).WithField("notification_id", notificationId).Error("Failed to store notification")
			return err
		}
	}

	// Update the status in storage
	if err := nm.storage.UpdateNotificationStatus(notificationId, notificationStatus, ""); err != nil {
		logrus.WithError(err).WithField("notification_id", notificationId).Error("Failed to update notification status")
		return err
	}

	// Log the status change with old status if available
	oldStatus := "unknown"
	if existingRecord != nil {
		oldStatus = string(existingRecord.Status)
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": notificationId,
		"old_status":      oldStatus,
		"new_status":      status,
	}).Info("Notification status updated successfully")

	return nil
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

// GetTemplateByID returns a template by ID (latest version)
func (nm *NotificationManagerImpl) GetTemplateByID(templateID string) (*models.Template, error) {
	return nm.templateManager.GetTemplateByID(templateID)
}

// GetTemplateByIDAndVersion returns a specific version of a template
func (nm *NotificationManagerImpl) GetTemplateByIDAndVersion(templateID string, version int) (*models.Template, error) {
	return nm.templateManager.GetTemplateByIDAndVersion(templateID, version)
}
