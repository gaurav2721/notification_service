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

	// Parse the payload to extract notification details
	var notificationData map[string]interface{}
	if err := json.Unmarshal([]byte(message.Payload), &notificationData); err != nil {
		logrus.WithError(err).Error("Failed to parse notification payload")
		return fmt.Errorf("failed to parse notification payload: %w", err)
	}

	// Extract recipient information
	recipientData, ok := notificationData["recipient"].(map[string]interface{})
	if !ok {
		logrus.Error("Invalid recipient data in notification payload")
		return fmt.Errorf("invalid recipient data in notification payload")
	}

	// Extract email address from recipient
	email, ok := recipientData["email"].(string)
	if !ok || email == "" {
		logrus.Error("Missing or invalid email address in recipient data")
		return fmt.Errorf("missing or invalid email address in recipient data")
	}

	// Extract content from notification
	content, ok := notificationData["content"].(map[string]interface{})
	if !ok {
		logrus.Error("Invalid content data in notification payload")
		return fmt.Errorf("invalid content data in notification payload")
	}

	// Extract "from" information for email
	var fromEmail string
	if fromData, ok := notificationData["from"].(map[string]interface{}); ok {
		if from, ok := fromData["email"].(string); ok {
			fromEmail = from
		}
	}

	// If no "from" email in notification, we'll use the default from the email service
	if fromEmail == "" {
		logrus.Debug("No 'from' email specified, will use default from email service")
	}

	// Create EmailNotificationRequest
	emailNotification := &models.EmailNotificationRequest{
		ID:   message.ID,
		Type: string(message.Type),
		Content: models.EmailContent{
			Subject:   getStringFromMap(content, "subject"),
			EmailBody: getStringFromMap(content, "email_body"),
		},
		Recipients: []string{email},
	}

	// Add "from" field if available
	if fromEmail != "" {
		emailNotification.From = &models.EmailSender{
			Email: fromEmail,
		}
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"to":              email,
		"from":            fromEmail,
		"subject":         emailNotification.Content.Subject,
	}).Info("Sending email notification")

	// Send email using the email service
	response, err := ep.emailService.SendEmail(ctx, emailNotification)
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
