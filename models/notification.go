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

// NotificationPriority represents the priority level
type NotificationPriority string

const (
	LowPriority    NotificationPriority = "low"
	NormalPriority NotificationPriority = "normal"
	HighPriority   NotificationPriority = "high"
	UrgentPriority NotificationPriority = "urgent"
)

// Notification represents a notification request
type Notification struct {
	ID          string                 `json:"id"`
	Type        NotificationType       `json:"type"`
	Priority    NotificationPriority   `json:"priority"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	TemplateID  string                 `json:"template_id,omitempty"`
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
func NewNotification(notificationType NotificationType, title, message string, recipients []string) *Notification {
	return &Notification{
		ID:         uuid.New().String(),
		Type:       notificationType,
		Priority:   NormalPriority,
		Title:      title,
		Message:    message,
		Recipients: recipients,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  time.Now(),
		Status:     "pending",
	}
}

// NewTemplate creates a new notification template
func NewTemplate(name string, notificationType NotificationType, subject, body string, variables []string) *NotificationTemplate {
	now := time.Now()
	return &NotificationTemplate{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      notificationType,
		Subject:   subject,
		Body:      body,
		Variables: variables,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
