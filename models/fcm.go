package models

import (
	"fmt"
	"strings"
	"time"
)

// FCMNotificationRequest represents an FCM notification request
type FCMNotificationRequest struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Content    FCMContent `json:"content"`
	Recipients []string   `json:"recipients"`
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

	// Validate recipients
	if len(notification.Recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	for i, recipient := range notification.Recipients {
		if recipient == "" {
			return fmt.Errorf("recipient at index %d cannot be empty", i)
		}
		if err := ValidateFCMDeviceToken(recipient); err != nil {
			return fmt.Errorf("invalid device token at index %d: %s", i, recipient)
		}
	}

	return nil
}

// ValidateFCMDeviceToken validates if a string is a valid FCM device token
func ValidateFCMDeviceToken(deviceToken string) error {
	if deviceToken == "" {
		return fmt.Errorf("device token cannot be empty")
	}

	deviceToken = strings.TrimSpace(deviceToken)

	// FCM device tokens (Firebase Instance ID tokens) are typically 140+ characters long
	// and contain alphanumeric characters, hyphens, and underscores
	if len(deviceToken) < 140 {
		return fmt.Errorf("FCM device token must be at least 140 characters long")
	}

	// Check if the device token contains only valid characters
	for _, char := range deviceToken {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			char == '-' || char == '_' || char == ':') {
			return fmt.Errorf("FCM device token must contain only alphanumeric characters, hyphens, underscores, and colons")
		}
	}

	return nil
}
