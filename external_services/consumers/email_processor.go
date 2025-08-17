package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gaurav2721/notification-service/external_services/email"
	"github.com/gaurav2721/notification-service/models"
	"github.com/sirupsen/logrus"
)

// emailProcessor handles email notification processing
type emailProcessor struct {
	emailService email.EmailService
}

// NewEmailProcessor creates a new email processor
func NewEmailProcessor() NotificationProcessor {
	// Create email service directly
	emailService := email.NewEmailService()

	return &emailProcessor{
		emailService: emailService,
	}
}

// NewEmailProcessorWithService creates a new email processor with a specific email service
func NewEmailProcessorWithService(emailService email.EmailService) NotificationProcessor {
	return &emailProcessor{
		emailService: emailService,
	}
}

// ProcessNotification processes an email notification
func (ep *emailProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
		"timestamp":       message.Timestamp,
	}).Debug("Processing email notification")

	// Parse the payload directly into EmailNotificationRequest
	var emailNotification models.EmailNotificationRequest
	if err := json.Unmarshal([]byte(message.Payload), &emailNotification); err != nil {
		logrus.WithError(err).Error("Failed to parse notification payload into EmailNotificationRequest")
		return fmt.Errorf("failed to parse notification payload into EmailNotificationRequest: %w", err)
	}

	// Validate the parsed notification
	if emailNotification.Recipient == "" {
		logrus.Error("No recipient specified in email notification")
		return fmt.Errorf("no recipient specified in email notification")
	}

	// Use the message ID if not set in the notification
	if emailNotification.ID == "" {
		emailNotification.ID = message.ID
	}

	// Use the message type if not set in the notification
	if emailNotification.Type == "" {
		emailNotification.Type = string(message.Type)
	}

	// Extract recipient and from information for logging
	recipient := emailNotification.Recipient

	fromEmail := ""
	if emailNotification.From != nil {
		fromEmail = emailNotification.From.Email
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"to":              recipient,
		"from":            fromEmail,
		"subject":         emailNotification.Content.Subject,
	}).Info("Sending email notification")

	// Send email using the email service
	response, err := ep.emailService.SendEmail(ctx, &emailNotification)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"notification_id": message.ID,
			"error":           err.Error(),
		}).Error("Failed to send email notification")
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log successful email sending
	if emailResponse, ok := response.(*models.EmailResponse); ok {
		logrus.WithFields(logrus.Fields{
			"notification_id": message.ID,
			"status":          emailResponse.Status,
			"sent_at":         emailResponse.SentAt,
		}).Info("Email notification sent successfully")
	} else {
		logrus.WithField("notification_id", message.ID).Info("Email notification sent successfully")
	}

	return nil
}

// GetNotificationType returns the notification type this processor handles
func (ep *emailProcessor) GetNotificationType() NotificationType {
	return EmailNotification
}
