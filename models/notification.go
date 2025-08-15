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
