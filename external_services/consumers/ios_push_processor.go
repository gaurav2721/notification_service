package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gaurav2721/notification-service/models"
	"github.com/sirupsen/logrus"
)

// iosPushProcessor handles iOS push notification processing
type iosPushProcessor struct {
	apnsService interface {
		SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	}
}

// NewIOSPushProcessor creates a new iOS push notification processor
func NewIOSPushProcessor() NotificationProcessor {
	return &iosPushProcessor{}
}

// NewIOSPushProcessorWithService creates a new iOS push processor with a specific APNS service
func NewIOSPushProcessorWithService(apnsService interface {
	SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
	return &iosPushProcessor{
		apnsService: apnsService,
	}
}

// ProcessNotification processes an iOS push notification
func (ip *iosPushProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
		"timestamp":       message.Timestamp,
	}).Debug("Processing iOS push notification")

	// If no APNS service is available, just log and return
	if ip.apnsService == nil {
		logrus.Warn("No APNS service available, skipping iOS push notification")
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

	// Extract device tokens from recipients
	recipients, ok := notificationData["recipients"].([]interface{})
	if !ok {
		logrus.Error("Invalid recipients data in notification payload")
		return fmt.Errorf("invalid recipients data in notification payload")
	}

	// Convert recipients to string slice
	deviceTokens := make([]string, 0, len(recipients))
	for _, recipient := range recipients {
		if token, ok := recipient.(string); ok {
			deviceTokens = append(deviceTokens, token)
		}
	}

	// Create APNS notification request
	apnsNotification := &models.APNSNotificationRequest{
		ID:   message.ID,
		Type: string(message.Type),
		Content: models.APNSContent{
			Title: getStringFromMap(content, "title"),
			Body:  getStringFromMap(content, "body"),
		},
		Recipients: deviceTokens,
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"device_tokens":   len(deviceTokens),
		"title":           apnsNotification.Content.Title,
		"body":            apnsNotification.Content.Body,
	}).Info("Sending iOS push notification")

	// Send push notification using the APNS service
	response, err := ip.apnsService.SendPushNotification(ctx, apnsNotification)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"notification_id": message.ID,
			"error":           err.Error(),
		}).Error("Failed to send iOS push notification")
		return fmt.Errorf("failed to send iOS push notification: %w", err)
	}

	// Log successful push notification sending
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"response":        response,
	}).Info("iOS push notification sent successfully")

	return nil
}

// GetNotificationType returns the notification type this processor handles
func (ip *iosPushProcessor) GetNotificationType() NotificationType {
	return IOSPushNotification
}
