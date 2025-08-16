package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// EmailNotificationRequest represents an email notification request
type EmailNotificationRequest struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Content    EmailContent `json:"content"`
	Recipients []string     `json:"recipients"`
	From       *EmailSender `json:"from,omitempty"`
}

// EmailContent represents the content of an email notification
type EmailContent struct {
	Subject   string `json:"subject"`
	EmailBody string `json:"email_body"`
}

// EmailSender represents the sender information
type EmailSender struct {
	Email string `json:"email"`
}

// EmailResponse represents the response from email sending
type EmailResponse struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
	Channel string    `json:"channel"`
}

// ValidateEmailNotification validates the email notification request
func ValidateEmailNotification(notification *EmailNotificationRequest) error {
	if notification == nil {
		return fmt.Errorf("email notification cannot be nil")
	}

	if notification.ID == "" {
		return fmt.Errorf("email notification ID is required")
	}

	if notification.Type == "" {
		return fmt.Errorf("email notification type is required")
	}

	// Validate content
	if notification.Content.Subject == "" {
		return fmt.Errorf("email subject is required")
	}

	if notification.Content.EmailBody == "" {
		return fmt.Errorf("email body is required")
	}

	// Validate recipients
	if len(notification.Recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	for i, recipient := range notification.Recipients {
		if recipient == "" {
			return fmt.Errorf("recipient at index %d cannot be empty", i)
		}
		if _, err := mail.ParseAddress(recipient); err != nil {
			return fmt.Errorf("invalid email address at index %d: %s", i, recipient)
		}
	}

	// Validate from email if provided
	if notification.From != nil {
		if notification.From.Email == "" {
			return fmt.Errorf("from email cannot be empty if from field is provided")
		}
		if _, err := mail.ParseAddress(notification.From.Email); err != nil {
			return fmt.Errorf("invalid from email address: %s", notification.From.Email)
		}
	}

	return nil
}

// ValidateEmailAddress validates if a string is a valid email address
func ValidateEmailAddress(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}

	email = strings.TrimSpace(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}

	return nil
}
