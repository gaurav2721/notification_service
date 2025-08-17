package fcm

import (
	"context"
	"testing"

	"github.com/gaurav2721/notification-service/models"
)

func TestNewFCMService(t *testing.T) {
	service := NewFCMService()
	if service == nil {
		t.Fatal("Expected FCM service to be created, got nil")
	}
}

func TestSendPushNotification(t *testing.T) {
	service := NewFCMService()

	// Test with Android device token in recipient
	notification := &models.FCMNotificationRequest{
		ID:   "test_notification",
		Type: "android_push",
		Content: models.FCMContent{
			Title: "Test Title",
			Body:  "Test Body",
		},
		Recipient: "android_device_token_123",
	}

	// Test sending to Android device
	result, err := service.SendPushNotification(context.Background(), notification)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Test with no device token - should fail validation
	emptyNotification := &models.FCMNotificationRequest{
		ID:   "test_notification",
		Type: "android_push",
		Content: models.FCMContent{
			Title: "Test Title",
			Body:  "Test Body",
		},
		Recipient: "",
	}

	_, err = service.SendPushNotification(context.Background(), emptyNotification)
	if err == nil {
		t.Errorf("Expected error for empty recipient, got nil")
	}

	// Test with invalid notification
	_, err = service.SendPushNotification(context.Background(), "invalid")
	if err != ErrInvalidNotificationPayload {
		t.Errorf("Expected ErrInvalidNotificationPayload, got %v", err)
	}
}
