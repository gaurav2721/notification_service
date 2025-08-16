package models

import (
	"fmt"
	"strings"
	"time"
)

// APNSNotificationRequest represents an APNS notification request
type APNSNotificationRequest struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Content   APNSContent `json:"content"`
	Recipient string      `json:"recipient"`
}

// APNSContent represents the content of an APNS notification
type APNSContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// APNSSender represents the sender information for APNS
type APNSSender struct {
	DeviceToken string `json:"device_token"`
}

// APNSResponse represents the response from APNS push notification sending
type APNSResponse struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	SentAt       time.Time `json:"sent_at"`
	Channel      string    `json:"channel"`
	SuccessCount int       `json:"success_count"`
	FailureCount int       `json:"failure_count"`
}

// ValidateAPNSNotification validates the APNS notification request
func ValidateAPNSNotification(notification *APNSNotificationRequest) error {
	if notification == nil {
		return fmt.Errorf("APNS notification cannot be nil")
	}

	if notification.ID == "" {
		return fmt.Errorf("APNS notification ID is required")
	}

	if notification.Type == "" {
		return fmt.Errorf("APNS notification type is required")
	}

	// Validate content
	if notification.Content.Title == "" {
		return fmt.Errorf("APNS title is required")
	}

	if notification.Content.Body == "" {
		return fmt.Errorf("APNS body is required")
	}

	// Validate recipient
	if notification.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}

	if err := ValidateDeviceToken(notification.Recipient); err != nil {
		return fmt.Errorf("invalid device token: %s", err)
	}

	return nil
}

// ValidateDeviceToken validates if a string is a valid device token
func ValidateDeviceToken(deviceToken string) error {
	if deviceToken == "" {
		return fmt.Errorf("device token cannot be empty")
	}

	deviceToken = strings.TrimSpace(deviceToken)

	// iOS device tokens are typically 64 characters long and contain only hexadecimal characters
	if len(deviceToken) != 64 {
		return fmt.Errorf("device token must be exactly 64 characters long")
	}

	// Check if the device token contains only hexadecimal characters
	for _, char := range deviceToken {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return fmt.Errorf("device token must contain only hexadecimal characters (0-9, a-f, A-F)")
		}
	}

	return nil
}
