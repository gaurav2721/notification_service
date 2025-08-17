package models

import (
	"fmt"
	"strings"
	"time"
)

// SlackNotificationRequest represents a slack notification request
type SlackNotificationRequest struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	Content   SlackContent `json:"content"`
	Recipient string       `json:"recipient"`
}

// SlackContent represents the content of a slack notification
type SlackContent struct {
	Text string `json:"text"`
}

// SlackResponse represents the response from slack message sending
type SlackResponse struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
	Channel string    `json:"channel"`
}

// ValidateSlackNotification validates the slack notification request
func ValidateSlackNotification(notification *SlackNotificationRequest) error {
	if notification == nil {
		return fmt.Errorf("slack notification cannot be nil")
	}

	if notification.ID == "" {
		return fmt.Errorf("slack notification ID is required")
	}

	if notification.Type == "" {
		return fmt.Errorf("slack notification type is required")
	}

	// Validate content
	if notification.Content.Text == "" {
		return fmt.Errorf("slack text is required")
	}

	// Validate recipient
	if notification.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}

	if err := ValidateUserID(notification.Recipient); err != nil {
		return fmt.Errorf("invalid user ID: %s", err)
	}

	return nil
}

// ValidateUserID validates if a string is a valid user ID
func ValidateUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	userID = strings.TrimSpace(userID)
	if len(userID) < 1 {
		return fmt.Errorf("user ID must be at least 1 character long")
	}

	// Add any specific user ID validation rules here
	// For example, check for valid characters, length limits, etc.

	return nil
}
