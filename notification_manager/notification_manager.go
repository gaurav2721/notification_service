package notification_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gaurav2721/notification-service/external_services/kafka"
	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager/scheduler"
	"github.com/gaurav2721/notification-service/notification_manager/templates"
	"github.com/google/uuid"
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

// NewNotificationManagerWithDefaultTemplate creates a new notification manager with default template manager
// The scheduler is initialized internally within the notification manager
func NewNotificationManagerWithDefaultTemplate(
	userService interface{},
	kafkaService kafka.KafkaService,
) *NotificationManagerImpl {
	return &NotificationManagerImpl{
		userService:     userService,
		kafkaService:    kafkaService,
		scheduler:       scheduler.NewScheduler(),
		templateManager: templates.NewTemplateManager(),
		storage:         NewInMemoryStorage(),
	}
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
func (nm *NotificationManagerImpl) GetNotificationStatus(notificationID string) (interface{}, error) {
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
func (nm *NotificationManagerImpl) SetNotificationStatus(notificationId string, notification *models.NotificationRequest, status string) error {
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
func (nm *NotificationManagerImpl) CreateTemplate(template interface{}) (interface{}, error) {
	// Type assertion to get template
	tmpl, ok := template.(*models.Template)
	if !ok {
		return nil, ErrTemplateNotFound
	}

	return nm.templateManager.CreateTemplate(context.Background(), tmpl)
}

// GetTemplateVersion retrieves a specific version of a notification template
func (nm *NotificationManagerImpl) GetTemplateVersion(templateID string, version int) (interface{}, error) {
	return nm.templateManager.GetTemplateVersion(context.Background(), templateID, version)
}

// GetPredefinedTemplates returns all predefined templates
func (nm *NotificationManagerImpl) GetPredefinedTemplates() []*models.Template {
	return nm.templateManager.GetPredefinedTemplates()
}

// ProcessNotificationRequest handles the complete notification request processing
func (nm *NotificationManagerImpl) ProcessNotificationRequest(ctx context.Context, request *models.NotificationRequest) (interface{}, error) {
	logrus.Debug("Processing notification request")

	// Generate notification ID
	notificationID := nm.generateID()

	// Process template if provided and generate content
	if request.Template != nil {
		logrus.Debug("Processing template to generate content")
		generatedContent, err := nm.processTemplateToContent(request.Template, request.Type)
		if err != nil {
			logrus.WithError(err).Error("Failed to process template")
			return nil, fmt.Errorf("template processing failed: %v", err)
		}

		// Replace or merge the content with generated content
		if request.Content == nil {
			request.Content = generatedContent
		} else {
			// Merge generated content with existing content, giving priority to generated content
			for key, value := range generatedContent {
				request.Content[key] = value
			}
		}

		logrus.WithField("generated_content", generatedContent).Debug("Template content generated and merged")
	}

	// Check if it's a scheduled notification
	if request.ScheduledAt != nil {
		logrus.Debug("Processing scheduled notification")

		// Schedule notification with a job function
		err := nm.ScheduleNotification(ctx, notificationID, request, func() error {
			// This job will be executed at the scheduled time
			logrus.WithField("notification_id", notificationID).Info("Executing scheduled notification job")

			// Process notification for recipients
			_, err := nm.processNotificationForRecipients(ctx, request, notificationID)
			if err != nil {
				logrus.WithError(err).Error("Failed to process notification for recipients")
				// Set status to failed if processing fails
				if statusErr := nm.SetNotificationStatus(notificationID, request, "failed"); statusErr != nil {
					logrus.WithError(statusErr).WithField("notification_id", notificationID).Warn("Failed to set notification status to failed")
				}
				return err
			}

			// Set notification status to sent after successful processing
			if err := nm.SetNotificationStatus(notificationID, request, "sent"); err != nil {
				logrus.WithError(err).WithField("notification_id", notificationID).Warn("Failed to set notification status to sent")
			}

			return nil
		})

		if err != nil {
			logrus.WithError(err).Error("Failed to schedule notification")
			return nil, err
		}

		// Set notification status to scheduled
		if err := nm.SetNotificationStatus(notificationID, request, "scheduled"); err != nil {
			logrus.WithError(err).WithField("notification_id", notificationID).Warn("Failed to set notification status to scheduled")
		}

		logrus.WithField("notification_id", notificationID).Debug("Notification scheduled successfully")
		return map[string]interface{}{
			"id":     notificationID,
			"status": "scheduled",
		}, nil
	}

	// Process notification for recipients
	responses, err := nm.processNotificationForRecipients(ctx, request, notificationID)
	if err != nil {
		logrus.WithError(err).Error("Failed to process notification for recipients")
		// Set notification status to failed
		if statusErr := nm.SetNotificationStatus(notificationID, request, "failed"); statusErr != nil {
			logrus.WithError(statusErr).WithField("notification_id", notificationID).Warn("Failed to set notification status to failed")
		}
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"notification_id":  notificationID,
		"total_recipients": len(request.Recipients),
		"queued_count":     len(responses),
	}).Debug("Notification processing completed")

	// Set notification status to sent
	if err := nm.SetNotificationStatus(notificationID, request, "sent"); err != nil {
		logrus.WithError(err).WithField("notification_id", notificationID).Warn("Failed to set notification status to sent")
	}

	// Return aggregated response
	return map[string]interface{}{
		"id":     notificationID,
		"status": "sent",
	}, nil
}

// processTemplateToContent processes a template and returns the generated content
func (nm *NotificationManagerImpl) processTemplateToContent(template *models.TemplateData, notificationType string) (map[string]interface{}, error) {
	if template == nil {
		return nil, fmt.Errorf("template cannot be nil")
	}

	// Get the template from the template manager using the specified version
	templateObj, err := nm.templateManager.GetTemplateByIDAndVersion(template.ID, template.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %v", err)
	}

	// Validate that the template type matches the notification type
	if string(templateObj.Type) != notificationType {
		return nil, fmt.Errorf("template type %s does not match notification type %s", templateObj.Type, notificationType)
	}

	// Validate required variables
	if err := templateObj.ValidateRequiredVariables(template.Data); err != nil {
		return nil, fmt.Errorf("template validation failed: %v", err)
	}

	// Process template content based on type
	content := make(map[string]interface{})

	switch notificationType {
	case "email":
		// Process email template
		subject := nm.processTemplateString(templateObj.Content.Subject, template.Data)
		emailBody := nm.processTemplateString(templateObj.Content.EmailBody, template.Data)

		content["subject"] = subject
		content["email_body"] = emailBody

	case "slack":
		// Process slack template
		text := nm.processTemplateString(templateObj.Content.Text, template.Data)
		content["text"] = text

	case "in_app":
		// Process in-app template
		title := nm.processTemplateString(templateObj.Content.Title, template.Data)
		body := nm.processTemplateString(templateObj.Content.Body, template.Data)

		content["title"] = title
		content["body"] = body

	default:
		return nil, fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	logrus.WithFields(logrus.Fields{
		"template_id": template.ID,
		"type":        notificationType,
		"content":     content,
	}).Debug("Template processed successfully")

	return content, nil
}

// processTemplateString replaces template variables with actual values
func (nm *NotificationManagerImpl) processTemplateString(templateStr string, data map[string]interface{}) string {
	result := templateStr

	// Replace variables in the format {{variable_name}}
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		stringValue := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, stringValue)
	}

	return result
}

// processNotificationForRecipients processes notifications for all recipients
func (nm *NotificationManagerImpl) processNotificationForRecipients(ctx context.Context, request *models.NotificationRequest, notificationID string) ([]interface{}, error) {
	// Get recipient information from userService
	logrus.Debug("Fetching recipient information from user service")

	// Type assert userService to get the GetUsersByIDs method
	userService, ok := nm.userService.(interface {
		GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error)
		GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error)
	})
	if !ok {
		return nil, fmt.Errorf("userService does not implement required methods")
	}

	users, err := userService.GetUsersByIDs(ctx, request.Recipients)
	if err != nil {
		logrus.WithError(err).Error("Failed to get recipient information")
		return nil, fmt.Errorf("failed to get recipient information: %v", err)
	}

	logrus.WithField("valid_users", len(users)).Debug("Retrieved user information")

	if len(users) == 0 {
		logrus.Warn("No valid recipients found for notification")
		return nil, fmt.Errorf("no valid recipients found")
	}

	// Process notifications based on type and fetch relevant user information
	var responses []interface{}

	logrus.WithFields(logrus.Fields{
		"notification_id":   notificationID,
		"recipient_count":   len(users),
		"notification_type": request.Type,
	}).Debug("Processing notification for recipients")

	for _, user := range users {
		logrus.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Debug("Processing notification for user")

		// Get detailed user notification info based on notification type
		userNotificationInfo, err := userService.GetUserNotificationInfo(ctx, user.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID,
				"error":   err.Error(),
			}).Error("Failed to get user notification info")
			continue
		}

		// Process notification based on type
		userResponses, err := nm.processNotificationByType(notificationID, *request, userNotificationInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID,
				"error":   err.Error(),
			}).Error("Failed to process notification for user")
			continue
		}

		responses = append(responses, userResponses...)
	}

	return responses, nil
}

// processNotificationByType processes notifications based on type and user information
func (nm *NotificationManagerImpl) processNotificationByType(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo) ([]interface{}, error) {
	var responses []interface{}

	switch request.Type {
	case "email":
		// For email notifications, use email as recipient
		if userInfo.Email == "" {
			logrus.WithField("user_id", userInfo.ID).Warn("User has no email address")
			return responses, nil
		}

		// Create email-specific message
		emailMessage := nm.createEmailMessage(notificationID, request, userInfo)

		// Post to email channel
		err := nm.postToKafkaChannel("email", emailMessage)
		if err != nil {
			return responses, fmt.Errorf("failed to post email notification: %v", err)
		}

		// Create response
		response := &models.NotificationResponse{
			ID:      notificationID,
			Status:  "queued",
			Message: fmt.Sprintf("Email notification queued for user %s", userInfo.FullName),
			SentAt:  time.Now(),
			Channel: "email",
		}
		responses = append(responses, response)

	case "slack":
		// For slack notifications, use slack channel as recipient
		if userInfo.SlackChannel == "" {
			logrus.WithField("user_id", userInfo.ID).Warn("User has no slack channel")
			return responses, nil
		}

		// Create slack-specific message
		slackMessage := nm.createSlackMessage(notificationID, request, userInfo)

		// Post to slack channel
		err := nm.postToKafkaChannel("slack", slackMessage)
		if err != nil {
			return responses, fmt.Errorf("failed to post slack notification: %v", err)
		}

		// Create response
		response := &models.NotificationResponse{
			ID:      notificationID,
			Status:  "queued",
			Message: fmt.Sprintf("Slack notification queued for user %s", userInfo.FullName),
			SentAt:  time.Now(),
			Channel: "slack",
		}
		responses = append(responses, response)

	case "in_app":
		// For in_app notifications, determine push type based on user devices
		if len(userInfo.Devices) == 0 {
			logrus.WithField("user_id", userInfo.ID).Warn("User has no active devices")
			return responses, nil
		}

		// Group devices by type
		iosDevices := make([]*models.UserDeviceInfo, 0)
		androidDevices := make([]*models.UserDeviceInfo, 0)

		for _, device := range userInfo.Devices {
			if device.IsActive && device.DeviceToken != "" {
				switch device.DeviceType {
				case "ios":
					iosDevices = append(iosDevices, device)
				case "android":
					androidDevices = append(androidDevices, device)
				}
			}
		}

		// Send to iOS devices - one message per device token
		for _, device := range iosDevices {
			if device.IsActive && device.DeviceToken != "" {
				iosMessage := nm.createIndividualPushMessage(notificationID, request, userInfo, device.DeviceToken, "ios_push")
				err := nm.postToKafkaChannel("ios_push", iosMessage)
				if err != nil {
					logrus.WithError(err).WithField("device_token", device.DeviceToken).Error("Failed to post iOS push notification")
				} else {
					response := &models.NotificationResponse{
						ID:      notificationID,
						Status:  "queued",
						Message: fmt.Sprintf("iOS push notification queued for user %s (device: %s)", userInfo.FullName, device.DeviceToken[:8]+"..."),
						SentAt:  time.Now(),
						Channel: "ios_push",
					}
					responses = append(responses, response)
				}
			}
		}

		// Send to Android devices - one message per device token
		for _, device := range androidDevices {
			if device.IsActive && device.DeviceToken != "" {
				androidMessage := nm.createIndividualPushMessage(notificationID, request, userInfo, device.DeviceToken, "android_push")
				err := nm.postToKafkaChannel("android_push", androidMessage)
				if err != nil {
					logrus.WithError(err).WithField("device_token", device.DeviceToken).Error("Failed to post Android push notification")
				} else {
					response := &models.NotificationResponse{
						ID:      notificationID,
						Status:  "queued",
						Message: fmt.Sprintf("Android push notification queued for user %s (device: %s)", userInfo.FullName, device.DeviceToken[:8]+"..."),
						SentAt:  time.Now(),
						Channel: "android_push",
					}
					responses = append(responses, response)
				}
			}
		}
	default:
		return responses, fmt.Errorf("unsupported notification type: %s", request.Type)
	}

	return responses, nil
}

// postToKafkaChannel posts the notification message to the appropriate Kafka channel
func (nm *NotificationManagerImpl) postToKafkaChannel(notificationType string, message interface{}) error {
	// Convert message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification message: %v", err)
	}

	messageStr := string(messageJSON)

	// Post to appropriate channel based on notification type
	switch notificationType {
	case "email":
		select {
		case nm.kafkaService.GetEmailChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("email channel is full")
		}

	case "slack":
		select {
		case nm.kafkaService.GetSlackChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("slack channel is full")
		}

	case "ios_push":
		select {
		case nm.kafkaService.GetIOSPushNotificationChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("iOS push notification channel is full")
		}

	case "android_push":
		select {
		case nm.kafkaService.GetAndroidPushNotificationChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("android push notification channel is full")
		}

	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	return nil
}

// createEmailMessage creates an email-specific notification message
func (nm *NotificationManagerImpl) createEmailMessage(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo) *models.EmailNotificationRequest {
	// Extract content from request
	var subject, emailBody string
	if request.Content != nil {
		if subj, ok := request.Content["subject"].(string); ok {
			subject = subj
		}
		if body, ok := request.Content["email_body"].(string); ok {
			emailBody = body
		}
	}

	emailNotification := &models.EmailNotificationRequest{
		ID:   notificationID,
		Type: "email",
		Content: models.EmailContent{
			Subject:   subject,
			EmailBody: emailBody,
		},
		Recipient: userInfo.Email,
	}

	// Add from field if provided
	if request.From != nil {
		emailNotification.From = &models.EmailSender{
			Email: request.From.Email,
		}
	}

	return emailNotification
}

// createSlackMessage creates a slack-specific notification message
func (nm *NotificationManagerImpl) createSlackMessage(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo) *models.SlackNotificationRequest {
	// Extract content from request
	var text string
	if request.Content != nil {
		if txt, ok := request.Content["text"].(string); ok {
			text = txt
		}
	}

	slackNotification := &models.SlackNotificationRequest{
		ID:        notificationID,
		Type:      "slack",
		Content:   models.SlackContent{Text: text},
		Recipient: userInfo.SlackChannel,
	}

	return slackNotification
}

// createIndividualPushMessage creates a push notification message for a single device token
func (nm *NotificationManagerImpl) createIndividualPushMessage(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo, deviceToken string, pushType string) interface{} {
	// Extract content from request
	var title, body string
	if request.Content != nil {
		if t, ok := request.Content["title"].(string); ok {
			title = t
		}
		if b, ok := request.Content["body"].(string); ok {
			body = b
		}
	}

	// Return appropriate notification type based on pushType
	switch pushType {
	case "ios_push":
		return &models.APNSNotificationRequest{
			ID:        notificationID,
			Type:      "ios_push",
			Content:   models.APNSContent{Title: title, Body: body},
			Recipient: deviceToken,
		}
	case "android_push":
		return &models.FCMNotificationRequest{
			ID:        notificationID,
			Type:      "android_push",
			Content:   models.FCMContent{Title: title, Body: body},
			Recipient: deviceToken,
		}
	default:
		// Fallback to generic map for unsupported types
		return map[string]interface{}{
			"notification_id": notificationID,
			"type":            pushType,
			"content":         map[string]interface{}{"title": title, "body": body},
			"recipients":      []string{deviceToken},
		}
	}
}

// generateID generates a UUID for notification IDs
func (nm *NotificationManagerImpl) generateID() string {
	return uuid.New().String()
}
