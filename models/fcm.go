package models

import (
	"fmt"
	"time"
)

// FCMNotificationRequest represents an FCM notification request
type FCMNotificationRequest struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Content   FCMContent `json:"content"`
	Recipient string     `json:"recipient"`
}

// FCMContent represents the content of an FCM notification
type FCMContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// FCMSender represents the sender information for FCM
type FCMSender struct {
	DeviceToken string `json:"device_token"`
}

// FCMResponse represents the response from FCM push notification sending
type FCMResponse struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	SentAt       time.Time `json:"sent_at"`
	Channel      string    `json:"channel"`
	SuccessCount int       `json:"success_count"`
	FailureCount int       `json:"failure_count"`
}

// ValidateFCMNotification validates the FCM notification request
func ValidateFCMNotification(notification *FCMNotificationRequest) error {
	if notification == nil {
		return fmt.Errorf("FCM notification cannot be nil")
	}

	if notification.ID == "" {
		return fmt.Errorf("FCM notification ID is required")
	}

	if notification.Type == "" {
		return fmt.Errorf("FCM notification type is required")
	}

	// Validate content
	if notification.Content.Title == "" {
		return fmt.Errorf("FCM title is required")
	}

	if notification.Content.Body == "" {
		return fmt.Errorf("FCM body is required")
	}

	// Validate recipient
	if notification.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}

	// Device token validation - only check if it's not empty
	if notification.Recipient == "" {
		return fmt.Errorf("device token cannot be empty")
	}

	return nil
}
