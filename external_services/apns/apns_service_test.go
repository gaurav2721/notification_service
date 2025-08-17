package apns

import (
	"context"
	"testing"

	"github.com/gaurav2721/notification-service/models"
)

func TestNewAPNSService(t *testing.T) {
	service := NewAPNSService()
	if service == nil {
		t.Fatal("Expected APNS service to be created, got nil")
	}
}

func TestSendPushNotification(t *testing.T) {
	service := NewAPNSService()

	// Test with iOS device token in recipient
	notification := &models.APNSNotificationRequest{
		ID:   "test_notification",
		Type: "ios_push",
		Content: models.APNSContent{
			Title: "Test Title",
			Body:  "Test Body",
		},
		Recipient: "ios_device_token_123",
	}

	// Test sending to iOS device
	result, err := service.SendPushNotification(context.Background(), notification)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Test with no device token - should fail validation
	emptyNotification := &models.APNSNotificationRequest{
		ID:   "test_notification",
		Type: "ios_push",
		Content: models.APNSContent{
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
