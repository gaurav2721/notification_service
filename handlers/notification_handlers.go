package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService interface {
		SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
		ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
		GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
		CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
		GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
		GetPredefinedTemplates() []*models.Template
	}
	userService interface {
		GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error)
		GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error)
	}
	kafkaService interface {
		GetEmailChannel() chan string
		GetSlackChannel() chan string
		GetIOSPushNotificationChannel() chan string
		GetAndroidPushNotificationChannel() chan string
	}
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	notificationService interface {
		SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
		ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
		GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
		CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
		GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
		GetPredefinedTemplates() []*models.Template
	},
	userService interface {
		GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error)
		GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error)
	},
	kafkaService interface {
		GetEmailChannel() chan string
		GetSlackChannel() chan string
		GetIOSPushNotificationChannel() chan string
		GetAndroidPushNotificationChannel() chan string
	},
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		userService:         userService,
		kafkaService:        kafkaService,
	}
}

// SendNotification handles POST /notifications
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	logrus.Debug("Received notification send request")

	var request models.NotificationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Warn("Invalid request body for notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use comprehensive validation
	validator := validation.NewNotificationValidator()

	validationResult := validator.ValidateNotificationRequest(&request)
	if !validationResult.IsValid {
		logrus.WithField("errors", validationResult.Errors).Warn("Validation failed for notification request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationResult.Errors,
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"type":        request.Type,
		"recipients":  len(request.Recipients),
		"scheduled":   request.ScheduledAt != nil,
		"hasTemplate": request.Template != nil,
		"hasFrom":     request.From != nil,
	}).Debug("Processing notification request")

	// Check if it's a scheduled notification
	if request.ScheduledAt != nil {
		logrus.Debug("Processing scheduled notification")
		// Create notification object for scheduling
		notification := &struct {
			ID          string
			Type        string
			Content     map[string]interface{}
			Template    *models.TemplateData
			Recipients  []string
			ScheduledAt *time.Time
			From        *struct {
				Email string `json:"email"`
			}
		}{
			ID:          generateID(),
			Type:        request.Type,
			Content:     request.Content,
			Template:    request.Template,
			Recipients:  request.Recipients,
			ScheduledAt: request.ScheduledAt,
			From:        request.From,
		}

		// Schedule notification
		response, err := h.notificationService.ScheduleNotification(c.Request.Context(), notification)
		if err != nil {
			logrus.WithError(err).Error("Failed to schedule notification")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logrus.WithField("notification_id", notification.ID).Debug("Notification scheduled successfully")
		c.JSON(http.StatusOK, response)
		return
	}

	// Get recipient information from userService
	logrus.Debug("Fetching recipient information from user service")
	users, err := h.userService.GetUsersByIDs(c.Request.Context(), request.Recipients)
	if err != nil {
		logrus.WithError(err).Error("Failed to get recipient information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get recipient information: %v", err)})
		return
	}

	logrus.WithField("valid_users", len(users)).Debug("Retrieved user information")

	if len(users) == 0 {
		logrus.Warn("No valid recipients found for notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid recipients found"})
		return
	}

	// Process notifications based on type and fetch relevant user information
	var responses []interface{}
	notificationID := generateID()

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
		userResponses, err := h.processNotificationByType(notificationID, request, userNotificationInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID,
				"error":   err.Error(),
			}).Error("Failed to process notification for user")
			continue
		}

		responses = append(responses, userResponses...)
	}

	logrus.WithFields(logrus.Fields{
		"notification_id":  notificationID,
		"total_recipients": len(users),
		"queued_count":     len(responses),
	}).Debug("Notification processing completed")

	// Return aggregated response
	c.JSON(http.StatusOK, gin.H{
		"notification_id":  notificationID,
		"total_recipients": len(users),
		"queued_count":     len(responses),
		"responses":        responses,
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

	fmt.Println("---------------> gaurav messageStr", messageStr)

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
	var request models.TemplateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate template content
	if err := request.Content.ValidateTemplateContent(request.Type); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template ID is required"})
		return
	}

	versionStr := c.Param("version")
	if versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version is required"})
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
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

// Helper function to generate a simple ID
func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
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

	// Personalize content with user information
	if userInfo.FullName != "" {
		emailBody = strings.ReplaceAll(emailBody, "{{recipient_name}}", userInfo.FullName)
		emailBody = strings.ReplaceAll(emailBody, "{{name}}", userInfo.FullName)
	}
	if userInfo.Email != "" {
		emailBody = strings.ReplaceAll(emailBody, "{{recipient_email}}", userInfo.Email)
		emailBody = strings.ReplaceAll(emailBody, "{{email}}", userInfo.Email)
	}

	emailNotification := &models.EmailNotificationRequest{
		ID:   notificationID,
		Type: "email",
		Content: models.EmailContent{
			Subject:   subject,
			EmailBody: emailBody,
		},
		Recipients: []string{userInfo.Email},
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

	// Personalize content with user information
	if userInfo.FullName != "" {
		text = strings.ReplaceAll(text, "{{recipient_name}}", userInfo.FullName)
		text = strings.ReplaceAll(text, "{{name}}", userInfo.FullName)
	}
	if userInfo.SlackUserID != "" {
		text = strings.ReplaceAll(text, "{{recipient_slack_id}}", userInfo.SlackUserID)
		text = strings.ReplaceAll(text, "{{slack_id}}", userInfo.SlackUserID)
	}

	slackNotification := &models.SlackNotificationRequest{
		ID:         notificationID,
		Type:       "slack",
		Content:    models.SlackContent{Text: text},
		Recipients: []string{userInfo.SlackChannel},
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
			ID:         notificationID,
			Type:       "ios_push",
			Content:    models.APNSContent{Title: title, Body: body},
			Recipients: []string{deviceToken},
		}
	case "android_push":
		return &models.FCMNotificationRequest{
			ID:         notificationID,
			Type:       "android_push",
			Content:    models.FCMContent{Title: title, Body: body},
			Recipients: []string{deviceToken},
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
