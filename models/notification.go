package models

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	EmailNotification NotificationType = "email"
	SlackNotification NotificationType = "slack"
	InAppNotification NotificationType = "in_app"
)

// Content represents the content structure for different notification types
type Content struct {
	// For Slack notifications
	Text string `json:"text,omitempty"`

	// For In-App notifications
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`

	// For Email notifications
	Subject   string `json:"subject,omitempty"`
	EmailBody string `json:"email_body,omitempty"`
}

// TemplateData represents the data structure for template usage in notifications
type TemplateData struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

// TemplateContent represents the content structure for different template types
type TemplateContent struct {
	// For Email templates
	Subject   string `json:"subject,omitempty"`
	EmailBody string `json:"email_body,omitempty"`

	// For Slack templates
	Text string `json:"text,omitempty"`

	// For In-App templates
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

// Template represents a notification template with versioning
type Template struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Type              NotificationType `json:"type"`
	Version           int              `json:"version"`
	Content           TemplateContent  `json:"content"`
	RequiredVariables []string         `json:"required_variables"`
	Description       string           `json:"description,omitempty"`
	Status            string           `json:"status"`
	CreatedAt         time.Time        `json:"created_at"`
}

// TemplateRequest represents the request structure for creating templates
type TemplateRequest struct {
	Name              string           `json:"name" binding:"required"`
	Type              NotificationType `json:"type" binding:"required"`
	Content           TemplateContent  `json:"content" binding:"required"`
	RequiredVariables []string         `json:"required_variables" binding:"required"`
	Description       string           `json:"description,omitempty"`
}

// TemplateResponse represents the response structure for template operations
type TemplateResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Version   int       `json:"version"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TemplateVersion represents a specific version of a template
type TemplateVersion struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Type              NotificationType `json:"type"`
	Version           int              `json:"version"`
	Content           TemplateContent  `json:"content"`
	RequiredVariables []string         `json:"required_variables"`
	Description       string           `json:"description,omitempty"`
	Status            string           `json:"status"`
	CreatedAt         time.Time        `json:"created_at"`
}

// Notification represents a notification request
type Notification struct {
	ID          string                 `json:"id"`
	Type        NotificationType       `json:"type"`
	Content     Content                `json:"content"`
	Template    *TemplateData          `json:"template,omitempty"`
	Recipients  []string               `json:"recipients"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	Status      string                 `json:"status"`
}

// NotificationTemplate represents a reusable template
type NotificationTemplate struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Type      NotificationType `json:"type"`
	Subject   string           `json:"subject"`
	Body      string           `json:"body"`
	Variables []string         `json:"variables"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// NotificationResponse represents the response after sending a notification
type NotificationResponse struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
	Channel string    `json:"channel"`
}

// NewNotification creates a new notification with default values
func NewNotification(notificationType NotificationType, content Content, recipients []string) *Notification {
	return &Notification{
		ID:         uuid.New().String(),
		Type:       notificationType,
		Content:    content,
		Recipients: recipients,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  time.Now(),
		Status:     "pending",
	}
}

// NewNotificationWithTemplate creates a new notification with template
func NewNotificationWithTemplate(notificationType NotificationType, template *TemplateData, recipients []string) *Notification {
	return &Notification{
		ID:         uuid.New().String(),
		Type:       notificationType,
		Template:   template,
		Recipients: recipients,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  time.Now(),
		Status:     "pending",
	}
}

// NewTemplate creates a new notification template
func NewTemplate(name string, notificationType NotificationType, content TemplateContent, requiredVariables []string, description string) *Template {
	return &Template{
		ID:                uuid.New().String(),
		Name:              name,
		Type:              notificationType,
		Version:           1, // Will be incremented by the service
		Content:           content,
		RequiredVariables: requiredVariables,
		Description:       description,
		Status:            "created",
		CreatedAt:         time.Now(),
	}
}

// NewTemplateVersion creates a new version of an existing template
func NewTemplateVersion(templateID string, name string, notificationType NotificationType, content TemplateContent, requiredVariables []string, description string, version int) *TemplateVersion {
	return &TemplateVersion{
		ID:                templateID,
		Name:              name,
		Type:              notificationType,
		Version:           version,
		Content:           content,
		RequiredVariables: requiredVariables,
		Description:       description,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// ValidateTemplateContent validates that the content matches the template type
func (tc *TemplateContent) ValidateTemplateContent(templateType NotificationType) error {
	switch templateType {
	case EmailNotification:
		if tc.Subject == "" || tc.EmailBody == "" {
			return ErrInvalidTemplateContent
		}
	case SlackNotification:
		if tc.Text == "" {
			return ErrInvalidTemplateContent
		}
	case InAppNotification:
		if tc.Title == "" || tc.Body == "" {
			return ErrInvalidTemplateContent
		}
	default:
		return ErrInvalidTemplateType
	}
	return nil
}

// ValidateRequiredVariables checks if all required variables are provided
func (t *Template) ValidateRequiredVariables(data map[string]interface{}) error {
	for _, requiredVar := range t.RequiredVariables {
		if _, exists := data[requiredVar]; !exists {
			return ErrMissingRequiredVariable
		}
	}
	return nil
}

// ValidateRequiredVariables checks if all required variables are provided
func (t *TemplateVersion) ValidateRequiredVariables(data map[string]interface{}) error {
	for _, requiredVar := range t.RequiredVariables {
		if _, exists := data[requiredVar]; !exists {
			return ErrMissingRequiredVariable
		}
	}
	return nil
}
