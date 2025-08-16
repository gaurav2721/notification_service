package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gaurav2721/notification-service/models"
	"github.com/sirupsen/logrus"
)

// androidPushProcessor handles Android push notification processing
type androidPushProcessor struct {
	fcmService interface {
		SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	}
}

// NewAndroidPushProcessor creates a new Android push notification processor
func NewAndroidPushProcessor() NotificationProcessor {
	return &androidPushProcessor{}
}

// NewAndroidPushProcessorWithService creates a new Android push processor with a specific FCM service
func NewAndroidPushProcessorWithService(fcmService interface {
	SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
	return &androidPushProcessor{
		fcmService: fcmService,
	}
}

// ProcessNotification processes an Android push notification
func (ap *androidPushProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
		"timestamp":       message.Timestamp,
	}).Debug("Processing Android push notification")

	// If no FCM service is available, just log and return
	if ap.fcmService == nil {
		logrus.Warn("No FCM service available, skipping Android push notification")
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

	// Create FCM notification request
	fcmNotification := &models.FCMNotificationRequest{
		ID:   message.ID,
		Type: string(message.Type),
		Content: models.FCMContent{
			Title: getStringFromMap(content, "title"),
			Body:  getStringFromMap(content, "body"),
		},
		Recipients: deviceTokens,
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"device_tokens":   len(deviceTokens),
		"title":           fcmNotification.Content.Title,
		"body":            fcmNotification.Content.Body,
	}).Info("Sending Android push notification")

	// Send push notification using the FCM service
	response, err := ap.fcmService.SendPushNotification(ctx, fcmNotification)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"notification_id": message.ID,
			"error":           err.Error(),
		}).Error("Failed to send Android push notification")
		return fmt.Errorf("failed to send Android push notification: %w", err)
	}

	// Log successful push notification sending
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"response":        response,
	}).Info("Android push notification sent successfully")

	return nil
}

// GetNotificationType returns the notification type this processor handles
func (ap *androidPushProcessor) GetNotificationType() NotificationType {
	return AndroidPushNotification
}
