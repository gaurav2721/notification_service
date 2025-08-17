package models

import (
	"fmt"
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

	// Device token validation - only check if it's not empty
	if notification.Recipient == "" {
		return fmt.Errorf("device token cannot be empty")
	}

	return nil
}
