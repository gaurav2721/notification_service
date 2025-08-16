package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// EmailNotificationRequest represents an email notification request
type EmailNotificationRequest struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Content    map[string]interface{} `json:"content"`
	Recipients []string               `json:"recipients"`
	From       *EmailSender           `json:"from,omitempty"`
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
	// Validate ID
	if notification.ID == "" {
		return ErrMissingEmailID
	}

	// Validate Type
	if notification.Type == "" {
		return ErrMissingEmailType
	}

	// Validate Content
	if notification.Content == nil {
		return ErrMissingEmailContent
	}

	// Validate Recipients
	if notification.Recipients == nil || len(notification.Recipients) == 0 {
		return ErrEmptyEmailRecipients
	}

	// Validate each recipient email
	for i, recipient := range notification.Recipients {
		if recipient == "" {
			return fmt.Errorf("recipient at index %d is empty", i)
		}
		if !isValidEmail(recipient) {
			return fmt.Errorf("invalid email address at index %d: %s", i, recipient)
		}
	}

	// Validate From email
	if notification.From == nil || notification.From.Email == "" {
		return ErrMissingEmailFrom
	}
	if !isValidEmail(notification.From.Email) {
		return ErrInvalidEmail
	}

	// Validate required content fields
	subject, hasSubject := notification.Content["subject"]
	if !hasSubject {
		return ErrMissingEmailSubject
	}
	if subjectStr, ok := subject.(string); !ok || strings.TrimSpace(subjectStr) == "" {
		return ErrMissingEmailSubject
	}

	body, hasBody := notification.Content["email_body"]
	if !hasBody {
		return ErrMissingEmailBody
	}
	if bodyStr, ok := body.(string); !ok || strings.TrimSpace(bodyStr) == "" {
		return ErrMissingEmailBody
	}

	return nil
}

// isValidEmail validates email address format
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// NewEmailNotification creates a new email notification with validation
func NewEmailNotification(id, emailType string, content map[string]interface{}, recipients []string, fromEmail string) (*EmailNotificationRequest, error) {
	notification := &EmailNotificationRequest{
		ID:         id,
		Type:       emailType,
		Content:    content,
		Recipients: recipients,
		From: &EmailSender{
			Email: fromEmail,
		},
	}

	if err := ValidateEmailNotification(notification); err != nil {
		return nil, err
	}

	return notification, nil
}
