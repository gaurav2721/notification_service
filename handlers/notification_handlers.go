package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gaurav2721/notification-service/models"
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

	var request struct {
		Type        string                 `json:"type" binding:"required"`
		Content     map[string]interface{} `json:"content"`
		Template    *models.TemplateData   `json:"template,omitempty"`
		Recipients  []string               `json:"recipients" binding:"required"`
		ScheduledAt *time.Time             `json:"scheduled_at"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Warn("Invalid request body for notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logrus.WithFields(logrus.Fields{
		"type":        request.Type,
		"recipients":  len(request.Recipients),
		"scheduled":   request.ScheduledAt != nil,
		"hasTemplate": request.Template != nil,
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
		}{
			ID:          generateID(),
			Type:        request.Type,
			Content:     request.Content,
			Template:    request.Template,
			Recipients:  request.Recipients,
			ScheduledAt: request.ScheduledAt,
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

	// Create notification messages for each recipient and post to Kafka channels
	var responses []interface{}
	notificationID := generateID()

	logrus.WithFields(logrus.Fields{
		"notification_id": notificationID,
		"recipient_count": len(users),
	}).Debug("Processing notification for recipients")

	for _, user := range users {
		logrus.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Debug("Processing notification for user")

		// Create personalized notification message for this user
		notificationMessage := h.createNotificationMessage(notificationID, request, user)

		fmt.Println("---------------> gaurav notificationMessage", notificationMessage)

		// Post to relevant Kafka channel based on notification type
		err := h.postToKafkaChannel(request.Type, notificationMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID,
				"error":   err.Error(),
			}).Error("Failed to post notification to Kafka channel")
			// Log error but continue with other recipients
			logrus.WithFields(logrus.Fields{
				"user_id": user.ID,
				"error":   err.Error(),
			}).Error("Failed to post notification for user")
			continue
		}

		logrus.WithField("user_id", user.ID).Debug("Notification queued successfully for user")

		// Create response for this recipient
		response := &models.NotificationResponse{
			ID:      notificationID,
			Status:  "queued",
			Message: fmt.Sprintf("Notification queued for user %s", user.FullName),
			SentAt:  time.Now(),
			Channel: request.Type,
		}
		responses = append(responses, response)
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

// createNotificationMessage creates a personalized notification message for a specific user
func (h *NotificationHandler) createNotificationMessage(notificationID string, request struct {
	Type        string                 `json:"type" binding:"required"`
	Content     map[string]interface{} `json:"content"`
	Template    *models.TemplateData   `json:"template,omitempty"`
	Recipients  []string               `json:"recipients" binding:"required"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
}, user *models.User) map[string]interface{} {
	// Create base notification message
	message := map[string]interface{}{
		"notification_id": notificationID,
		"type":            request.Type,
		"content":         request.Content,
		"template":        request.Template,
		"created_at":      time.Now(),
		"recipient": map[string]interface{}{
			"user_id":       user.ID,
			"email":         user.Email,
			"full_name":     user.FullName,
			"slack_user_id": user.SlackUserID,
			"slack_channel": user.SlackChannel,
			"phone_number":  user.PhoneNumber,
		},
	}

	// Add personalized content based on user information
	if request.Content != nil {
		// Personalize content with user information
		personalizedContent := make(map[string]interface{})
		for key, value := range request.Content {
			personalizedContent[key] = value
		}

		// Add user-specific personalization
		if user.FullName != "" {
			personalizedContent["recipient_name"] = user.FullName
		}
		if user.Email != "" {
			personalizedContent["recipient_email"] = user.Email
		}

		message["content"] = personalizedContent
	}

	return message
}

// postToKafkaChannel posts the notification message to the appropriate Kafka channel
func (h *NotificationHandler) postToKafkaChannel(notificationType string, message map[string]interface{}) error {
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
			return fmt.Errorf("Android push notification channel is full")
		}

	case "in_app":
		// For in-app notifications, we can use either iOS or Android channel
		// or create a separate in-app channel. For now, using iOS channel.
		select {
		case h.kafkaService.GetIOSPushNotificationChannel() <- messageStr:
			// Message sent successfully
		default:
			return fmt.Errorf("in-app notification channel is full")
		}

	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	return nil
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
