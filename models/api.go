package models

import (
	"time"
)

// NotificationRequest represents the request structure for sending notifications
type NotificationRequest struct {
	Type        string                 `json:"type" binding:"required"`
	Content     map[string]interface{} `json:"content"`
	Template    *TemplateData          `json:"template,omitempty"`
	Recipients  []string               `json:"recipients" binding:"required"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
	From        *struct {
		Email string `json:"email"`
	} `json:"from,omitempty"`
}
