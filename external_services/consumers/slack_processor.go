package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gaurav2721/notification-service/models"
	"github.com/sirupsen/logrus"
)

// slackProcessor handles slack notification processing
type slackProcessor struct {
	slackService interface {
		SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
	}
}

// NewSlackProcessor creates a new slack processor
func NewSlackProcessor() NotificationProcessor {
	return &slackProcessor{}
}

// NewSlackProcessorWithService creates a new slack processor with a specific slack service
func NewSlackProcessorWithService(slackService interface {
	SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
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

	// Parse the payload to extract notification details
	var notificationData map[string]interface{}
	if err := json.Unmarshal([]byte(message.Payload), &notificationData); err != nil {
		logrus.WithError(err).Error("Failed to parse notification payload")
		return fmt.Errorf("failed to parse notification payload: %w", err)
	}

	// Extract content from notification
	content, ok := notificationData["content"].(map[string]interface{})
	if !ok {
		logrus.Error("Invalid content data in notification payload")
		return fmt.Errorf("invalid content data in notification payload")
	}

	// Extract channel from notification (optional)
	var channel string
	if channelData, ok := notificationData["channel"].(string); ok {
		channel = channelData
	}

	// Create slack notification request
	slackNotification := &models.SlackNotificationRequest{
		ID:   message.ID,
		Type: string(message.Type),
		Content: models.SlackContent{
			Text: getStringFromMap(content, "text"),
		},
		Recipients: []string{}, // Slack doesn't use individual recipients like email
	}

	// Add channel if available
	if channel != "" {
		// Note: Channel is not part of the standard model, but could be added as metadata
		// For now, we'll include it in the content or handle it separately
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"channel":         channel,
		"text":            slackNotification.Content.Text,
	}).Info("Sending slack notification")

	// Send slack message using the slack service
	response, err := sp.slackService.SendSlackMessage(ctx, slackNotification)
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
