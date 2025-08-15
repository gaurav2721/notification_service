package email

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// EmailServiceImpl implements the EmailService interface
type EmailServiceImpl struct {
	dialer *gomail.Dialer
}

// NewEmailService creates a new email service instance
func NewEmailService() EmailService {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 587 // default SMTP port
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return &EmailServiceImpl{
		dialer: dialer,
	}
}

// NewEmailServiceWithConfig creates a new email service with custom configuration
func NewEmailServiceWithConfig(config *EmailConfig) EmailService {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

	return &EmailServiceImpl{
		dialer: dialer,
	}
}

// SendEmail sends an email notification
func (es *EmailServiceImpl) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*struct {
		ID         string
		Type       string
		Title      string
		Message    string
		Recipients []string
		Metadata   map[string]interface{}
	})
	if !ok {
		return nil, ErrEmailSendFailed
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USERNAME"))
	m.SetHeader("To", notif.Recipients...)
	m.SetHeader("Subject", notif.Title)
	m.SetBody("text/html", notif.Message)

	// Send email
	if err := es.dialer.DialAndSend(m); err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	// Return success response
	return &struct {
		ID      string    `json:"id"`
		Status  string    `json:"status"`
		Message string    `json:"message"`
		SentAt  time.Time `json:"sent_at"`
		Channel string    `json:"channel"`
	}{
		ID:      notif.ID,
		Status:  "sent",
		Message: "Email sent successfully",
		SentAt:  time.Now(),
		Channel: "email",
	}, nil
}
