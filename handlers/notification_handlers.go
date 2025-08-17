package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gaurav2721/notification-service/external_services/kafka"
	"github.com/gaurav2721/notification-service/external_services/user"
	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService notification_manager.NotificationManager
	userService         user.UserService
	kafkaService        kafka.KafkaService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	notificationService notification_manager.NotificationManager,
	userService user.UserService,
	kafkaService kafka.KafkaService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		userService:         userService,
		kafkaService:        kafkaService,
	}
}

// processTemplateToContent processes a template and returns the generated content
func (h *NotificationHandler) processTemplateToContent(template *models.TemplateData, notificationType string) (map[string]interface{}, error) {
	if template == nil {
		return nil, fmt.Errorf("template cannot be nil")
	}

	// Get the template from the notification service using the specified version
	templateObj, err := h.notificationService.GetTemplateByIDAndVersion(template.ID, template.Version)

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
		subject := h.processTemplateString(templateObj.Content.Subject, template.Data)
		emailBody := h.processTemplateString(templateObj.Content.EmailBody, template.Data)

		content["subject"] = subject
		content["email_body"] = emailBody

	case "slack":
		// Process slack template
		text := h.processTemplateString(templateObj.Content.Text, template.Data)
		content["text"] = text

	case "in_app":
		// Process in-app template
		title := h.processTemplateString(templateObj.Content.Title, template.Data)
		body := h.processTemplateString(templateObj.Content.Body, template.Data)

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
func (h *NotificationHandler) processTemplateString(templateStr string, data map[string]interface{}) string {
	result := templateStr

	// Replace variables in the format {{variable_name}}
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		stringValue := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, stringValue)
	}

	return result
}

// SendNotification handles POST /notifications
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	logrus.Debug("Received notification send request")

	// Get validated request from middleware
	validatedRequestInterface, exists := c.Get("validated_request")
	if !exists {
		logrus.Error("Validated request not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	requestPtr, ok := validatedRequestInterface.(*models.NotificationRequest)
	if !ok {
		logrus.Error("Failed to cast validated request to NotificationRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Dereference the pointer to get the actual request
	request := *requestPtr

	logrus.WithFields(logrus.Fields{
		"type":        request.Type,
		"recipients":  len(request.Recipients),
		"scheduled":   request.ScheduledAt != nil,
		"hasTemplate": request.Template != nil,
		"hasFrom":     request.From != nil,
	}).Debug("Processing notification request")

	// Process template if provided and generate content
	if request.Template != nil {
		logrus.Debug("Processing template to generate content")
		generatedContent, err := h.processTemplateToContent(request.Template, request.Type)
		if err != nil {
			logrus.WithError(err).Error("Failed to process template")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Template processing failed: %v", err)})
			return
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

	id := generateID()
	// Check if it's a scheduled notification
	if request.ScheduledAt != nil {
		logrus.Debug("Processing scheduled notification")
		// Create notification object for scheduling

		// Schedule notification with a job function
		err := h.notificationService.ScheduleNotification(c.Request.Context(), id, &request, func() error {
			// This job will be executed at the scheduled time
			// It should send the notification
			logrus.WithField("notification_id", id).Info("Executing scheduled notification job")

			// Get recipient information from userService
			_, err := h.processNotificationForRecipients(c, &request, id)
			if err != nil {
				logrus.WithError(err).Error("Failed to process notification for recipients")
				// Set status to failed if processing fails
				if statusErr := h.notificationService.SetNotificationStatus(c.Request.Context(), id, &request, "failed"); statusErr != nil {
					logrus.WithError(statusErr).WithField("notification_id", id).Warn("Failed to set notification status to failed")
				}
				return err
			}

			// Set notification status to sent after successful processing
			if err := h.notificationService.SetNotificationStatus(c.Request.Context(), id, &request, "sent"); err != nil {
				logrus.WithError(err).WithField("notification_id", id).Warn("Failed to set notification status to sent")
				// Continue even if status update fails
			}

			// TODO: Implement the actual notification sending logic here
			// For now, just log that the job would be executed
			return nil
		})

		if err != nil {
			logrus.WithError(err).Error("Failed to schedule notification")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Set notification status to sent after successful processing
		if err := h.notificationService.SetNotificationStatus(c.Request.Context(), id, &request, "scheduled"); err != nil {
			logrus.WithError(err).WithField("notification_id", id).Warn("Failed to set notification status to sent")
			// Continue even if status update fails
		}

		logrus.WithField("notification_id", id).Debug("Notification scheduled successfully")
		c.JSON(http.StatusOK, gin.H{
			"id":     id,
			"status": "scheduled",
		})
		return
	}

	// Get recipient information from userService
	responses, err := h.processNotificationForRecipients(c, &request, id)
	if err != nil {
		logrus.WithError(err).Error("Failed to process notification for recipients")
		// Set notification status to failed
		if statusErr := h.notificationService.SetNotificationStatus(c.Request.Context(), id, &request, "failed"); statusErr != nil {
			logrus.WithError(statusErr).WithField("notification_id", id).Warn("Failed to set notification status to failed")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.WithFields(logrus.Fields{
		"notification_id":  id,
		"total_recipients": len(request.Recipients),
		"queued_count":     len(responses),
	}).Debug("Notification processing completed")

	// Set notification status to sent
	if err := h.notificationService.SetNotificationStatus(c.Request.Context(), id, &request, "sent"); err != nil {
		logrus.WithError(err).WithField("notification_id", id).Warn("Failed to set notification status to sent")
		// Continue with response even if status update fails
	}

	// Return aggregated response
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"status": "sent",
	})
}

// postToKafkaChannel posts the notification message to the appropriate Kafka channel
func (h *NotificationHandler) postToKafkaChannel(notificationType string, message interface{}) error {
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
		case h.kafkaService.GetEmailChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("email channel is full")
		}

	case "slack":
		select {
		case h.kafkaService.GetSlackChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("slack channel is full")
		}

	case "ios_push":
		select {
		case h.kafkaService.GetIOSPushNotificationChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("iOS push notification channel is full")
		}

	case "android_push":
		select {
		case h.kafkaService.GetAndroidPushNotificationChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("android push notification channel is full")
		}

	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	return nil
}

// processNotificationByType processes notifications based on type and user information
func (h *NotificationHandler) processNotificationByType(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo) ([]interface{}, error) {
	var responses []interface{}

	switch request.Type {
	case "email":
		// For email notifications, use email as recipient
		if userInfo.Email == "" {
			logrus.WithField("user_id", userInfo.ID).Warn("User has no email address")
			return responses, nil
		}

		// Create email-specific message
		emailMessage := h.createEmailMessage(notificationID, request, userInfo)

		// Post to email channel
		err := h.postToKafkaChannel("email", emailMessage)
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
		slackMessage := h.createSlackMessage(notificationID, request, userInfo)

		// Post to slack channel
		err := h.postToKafkaChannel("slack", slackMessage)
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
				iosMessage := h.createIndividualPushMessage(notificationID, request, userInfo, device.DeviceToken, "ios_push")
				err := h.postToKafkaChannel("ios_push", iosMessage)
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
				androidMessage := h.createIndividualPushMessage(notificationID, request, userInfo, device.DeviceToken, "android_push")
				err := h.postToKafkaChannel("android_push", androidMessage)
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

// processNotificationForRecipients processes notifications for all recipients
func (h *NotificationHandler) processNotificationForRecipients(c *gin.Context, request *models.NotificationRequest, notificationID string) ([]interface{}, error) {
	// Get recipient information from userService
	logrus.Debug("Fetching recipient information from user service")
	users, err := h.userService.GetUsersByIDs(c.Request.Context(), request.Recipients)
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
		userNotificationInfo, err := h.userService.GetUserNotificationInfo(c.Request.Context(), user.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID,
				"error":   err.Error(),
			}).Error("Failed to get user notification info")
			continue
		}

		// Process notification based on type
		userResponses, err := h.processNotificationByType(notificationID, *request, userNotificationInfo)
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

// GetNotificationStatus handles GET /notifications/:id
func (h *NotificationHandler) GetNotificationStatus(c *gin.Context) {
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification ID is required"})
		return
	}

	response, err := h.notificationService.GetNotificationStatus(c.Request.Context(), notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateTemplate handles POST /templates
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	// Get validated request from middleware
	validatedRequestInterface, exists := c.Get("validated_template_request")
	if !exists {
		logrus.Error("Validated template request not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	requestPtr, ok := validatedRequestInterface.(*models.TemplateRequest)
	if !ok {
		logrus.Error("Failed to cast validated request to TemplateRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Dereference the pointer to get the actual request
	request := *requestPtr

	// Create template object
	template := &models.Template{
		Name:              request.Name,
		Type:              request.Type,
		Content:           request.Content,
		RequiredVariables: request.RequiredVariables,
		Description:       request.Description,
	}

	response, err := h.notificationService.CreateTemplate(c.Request.Context(), template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPredefinedTemplates handles GET /templates/predefined
func (h *NotificationHandler) GetPredefinedTemplates(c *gin.Context) {
	templates := h.notificationService.GetPredefinedTemplates()

	// Convert to response format
	var response []map[string]interface{}
	for _, template := range templates {
		response = append(response, map[string]interface{}{
			"id":                 template.ID,
			"name":               template.Name,
			"type":               string(template.Type),
			"version":            template.Version,
			"content":            template.Content,
			"description":        template.Description,
			"required_variables": template.RequiredVariables,
			"status":             template.Status,
			"created_at":         template.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": response,
		"count":     len(response),
	})
}

// GetTemplateVersion handles GET /templates/:templateId/versions/:version
func (h *NotificationHandler) GetTemplateVersion(c *gin.Context) {
	templateID := c.Param("templateId")
	versionStr := c.Param("version")

	// Parameters are already validated by middleware
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		// This should not happen as middleware validates this
		logrus.WithError(err).Error("Failed to parse version parameter")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	template, err := h.notificationService.GetTemplateVersion(c.Request.Context(), templateID, version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// HealthCheck handles GET /health
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "notification-service",
	})
}

// Helper function to generate a UUID
func generateID() string {
	return uuid.New().String()
}

// createEmailMessage creates an email-specific notification message
func (h *NotificationHandler) createEmailMessage(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo) *models.EmailNotificationRequest {
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
func (h *NotificationHandler) createSlackMessage(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo) *models.SlackNotificationRequest {
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
func (h *NotificationHandler) createIndividualPushMessage(notificationID string, request models.NotificationRequest, userInfo *models.UserNotificationInfo, deviceToken string, pushType string) interface{} {
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
