package email

import (
	"context"
	"errors"
)

// EmailService interface defines methods for email notifications
type EmailService interface {
	SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// DefaultEmailConfig returns default email configuration
func DefaultEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     "localhost",
		SMTPPort:     587,
		SMTPUsername: "",
		SMTPPassword: "",
		FromEmail:    "noreply@company.com",
		FromName:     "Notification Service",
	}
}

// Email service errors
var (
	ErrEmailSendFailed       = errors.New("failed to send email")
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrEmailTemplateNotFound = errors.New("email template not found")
)
