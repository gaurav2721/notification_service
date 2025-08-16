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

// BulkNotificationResponse represents the response for bulk notification operations
type BulkNotificationResponse struct {
	NotificationID  string                  `json:"notification_id"`
	TotalRecipients int                     `json:"total_recipients"`
	QueuedCount     int                     `json:"queued_count"`
	Responses       []*NotificationResponse `json:"responses"`
}

// UserRequest represents the request structure for user operations
type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

// DeviceRegistrationRequest represents the request structure for device registration
type DeviceRegistrationRequest struct {
	DeviceToken string `json:"device_token" binding:"required"`
	DeviceType  string `json:"device_type" binding:"required"`
	AppVersion  string `json:"app_version"`
	OSVersion   string `json:"os_version"`
	DeviceModel string `json:"device_model"`
}

// DeviceUpdateRequest represents the request structure for device updates
type DeviceUpdateRequest struct {
	AppVersion  string `json:"app_version"`
	OSVersion   string `json:"os_version"`
	DeviceModel string `json:"device_model"`
}
