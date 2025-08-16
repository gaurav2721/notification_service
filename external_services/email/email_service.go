package email

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"gopkg.in/gomail.v2"
)

// EmailServiceImpl implements the EmailService interface
type EmailServiceImpl struct {
	dialer *gomail.Dialer
}

// NewEmailService creates a new email service instance
// It checks environment variables and returns mock service if config is incomplete
func NewEmailService() EmailService {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Check if all required environment variables are present and non-empty
	if host == "" || portStr == "" || username == "" || password == "" {
		return NewMockEmailService()
	}

	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587 // default SMTP port
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return &EmailServiceImpl{
		dialer: dialer,
	}
}

// SendEmail sends an email notification
func (es *EmailServiceImpl) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*models.EmailNotificationRequest)
	if !ok {
		return nil, ErrEmailSendFailed
	}

	// Validate the email notification
	if err := models.ValidateEmailNotification(notif); err != nil {
		return nil, fmt.Errorf("email validation failed: %w", err)
	}

	// Create email message
	m := gomail.NewMessage()

	// Use "from" field if provided, otherwise fall back to environment variable
	fromEmail := os.Getenv("SMTP_USERNAME")
	if notif.From != nil && notif.From.Email != "" {
		fromEmail = notif.From.Email
	}

	m.SetHeader("From", fromEmail)
	m.SetHeader("To", notif.Recipients...)

	// Extract subject and body from content
	subject := notif.Content.Subject
	body := notif.Content.EmailBody

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send email
	if err := es.dialer.DialAndSend(m); err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	// Return success response
	return &models.EmailResponse{
		ID:      notif.ID,
		Status:  "sent",
		Message: "Email sent successfully",
		SentAt:  time.Now(),
		Channel: "email",
	}, nil
}
