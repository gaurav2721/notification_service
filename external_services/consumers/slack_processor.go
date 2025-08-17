package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gaurav2721/notification-service/external_services/slack"
	"github.com/gaurav2721/notification-service/models"
	"github.com/sirupsen/logrus"
)

// slackProcessor handles slack notification processing
type slackProcessor struct {
	slackService slack.SlackService
}

// NewSlackProcessor creates a new slack processor
func NewSlackProcessor() NotificationProcessor {
	return &slackProcessor{}
}

// NewSlackProcessorWithService creates a new slack processor with a specific slack service
func NewSlackProcessorWithService(slackService slack.SlackService) NotificationProcessor {
	return &slackProcessor{
		slackService: slackService,
	}
}

// ProcessNotification processes a slack notification
func (sp *slackProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
		"timestamp":       message.Timestamp,
	}).Debug("Processing slack notification")

	// If no slack service is available, just log and return
	if sp.slackService == nil {
		logrus.Warn("No slack service available, skipping slack notification")
		return nil
	}

	// Parse the payload directly into SlackNotificationRequest
	var slackNotification models.SlackNotificationRequest
	if err := json.Unmarshal([]byte(message.Payload), &slackNotification); err != nil {
		logrus.WithError(err).Error("Failed to parse notification payload into SlackNotificationRequest")
		return fmt.Errorf("failed to parse notification payload into SlackNotificationRequest: %w", err)
	}

	// Use the message ID if not set in the notification
	if slackNotification.ID == "" {
		slackNotification.ID = message.ID
	}

	// Use the message type if not set in the notification
	if slackNotification.Type == "" {
		slackNotification.Type = string(message.Type)
	}

	// Extract channel information for logging
	channel := slackNotification.Recipient

	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"channel":         channel,
		"text":            slackNotification.Content.Text,
	}).Info("Sending slack notification")

	// Send slack message using the slack service
	response, err := sp.slackService.SendSlackMessage(ctx, &slackNotification)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"notification_id": message.ID,
			"error":           err.Error(),
		}).Error("Failed to send slack notification")
		return fmt.Errorf("failed to send slack message: %w", err)
	}

	// Log successful slack sending
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"response":        response,
	}).Info("Slack notification sent successfully")

	return nil
}

// GetNotificationType returns the notification type this processor handles
func (sp *slackProcessor) GetNotificationType() NotificationType {
	return SlackNotification
}
