package notification_manager

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/google/uuid"
)

// NotificationManagerImpl implements the NotificationManager interface
type NotificationManagerImpl struct {
	emailService  interface{}
	slackService  interface{}
	inappService  interface{}
	userService   interface{}
	scheduler     interface{}
	templates     map[string]*models.Template
	templateMutex sync.RWMutex
	initialized   bool
}

// NewNotificationManager creates a new notification manager instance
func NewNotificationManager(
	emailService interface{},
	slackService interface{},
	inappService interface{},
	userService interface{},
	scheduler interface{},
) *NotificationManagerImpl {
	nm := &NotificationManagerImpl{
		emailService:  emailService,
		slackService:  slackService,
		inappService:  inappService,
		userService:   userService,
		scheduler:     scheduler,
		templates:     make(map[string]*models.Template),
		templateMutex: sync.RWMutex{},
		initialized:   false,
	}

	// Load predefined templates on startup
	nm.loadPredefinedTemplates()

	return nm
}

// loadPredefinedTemplates loads all predefined templates into the manager
func (nm *NotificationManagerImpl) loadPredefinedTemplates() {
	nm.templateMutex.Lock()
	defer nm.templateMutex.Unlock()

	predefinedTemplates := models.PredefinedTemplates()

	for _, template := range predefinedTemplates {
		nm.templates[template.ID] = template
		log.Printf("Loaded predefined template: %s (ID: %s)", template.Name, template.ID)
	}

	nm.initialized = true
	log.Printf("Loaded %d predefined templates", len(predefinedTemplates))
}

// GetPredefinedTemplates returns all predefined templates
func (nm *NotificationManagerImpl) GetPredefinedTemplates() []*models.Template {
	nm.templateMutex.RLock()
	defer nm.templateMutex.RUnlock()

	var predefinedTemplates []*models.Template
	for _, template := range nm.templates {
		// Check if it's a predefined template by looking at the fixed IDs
		if isPredefinedTemplateID(template.ID) {
			predefinedTemplates = append(predefinedTemplates, template)
		}
	}

	return predefinedTemplates
}

// isPredefinedTemplateID checks if a template ID is one of the predefined ones
func isPredefinedTemplateID(templateID string) bool {
	predefinedIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000", // Welcome Email
		"550e8400-e29b-41d4-a716-446655440001", // Password Reset
		"550e8400-e29b-41d4-a716-446655440002", // Order Confirmation
		"550e8400-e29b-41d4-a716-446655440003", // System Alert
		"550e8400-e29b-41d4-a716-446655440004", // Deployment Notification
		"550e8400-e29b-41d4-a716-446655440005", // Order Status Update
		"550e8400-e29b-41d4-a716-446655440006", // Payment Reminder
	}

	for _, id := range predefinedIDs {
		if templateID == id {
			return true
		}
	}
	return false
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
		// Get template version
		template, err := nm.GetTemplateVersion(ctx, notif.Template.ID, 1) // Default to version 1 for now
		if err != nil {
			return nil, err
		}

		// Validate required variables
		templateVersion, ok := template.(*models.TemplateVersion)
		if !ok {
			return nil, ErrTemplateNotFound
		}

		if err := templateVersion.ValidateRequiredVariables(notif.Template.Data); err != nil {
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

	nm.templateMutex.Lock()
	defer nm.templateMutex.Unlock()

	// Generate ID if not provided
	if tmpl.ID == "" {
		tmpl.ID = uuid.New().String()
	}

	// Set version and status
	tmpl.Version = 1
	tmpl.Status = "created"
	tmpl.CreatedAt = time.Now()

	// Store template
	nm.templates[tmpl.ID] = tmpl

	// Return response
	return &models.TemplateResponse{
		ID:        tmpl.ID,
		Name:      tmpl.Name,
		Type:      string(tmpl.Type),
		Version:   tmpl.Version,
		Status:    tmpl.Status,
		CreatedAt: tmpl.CreatedAt,
	}, nil
}

// GetTemplateVersion retrieves a specific version of a notification template
func (nm *NotificationManagerImpl) GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error) {
	nm.templateMutex.RLock()
	defer nm.templateMutex.RUnlock()

	template, exists := nm.templates[templateID]
	if !exists {
		return nil, ErrTemplateNotFound
	}

	// For now, we only support version 1
	// In a real implementation, you would store multiple versions
	if version != 1 {
		return nil, ErrTemplateNotFound
	}

	// Return template version
	return &models.TemplateVersion{
		ID:                template.ID,
		Name:              template.Name,
		Type:              template.Type,
		Version:           template.Version,
		Content:           template.Content,
		RequiredVariables: template.RequiredVariables,
		Description:       template.Description,
		Status:            template.Status,
		CreatedAt:         template.CreatedAt,
	}, nil
}
