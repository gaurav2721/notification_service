package apns

import (
	"context"
	"testing"
	"time"
)

func TestNewAPNSService(t *testing.T) {
	service := NewAPNSService()
	if service == nil {
		t.Fatal("Expected APNS service to be created, got nil")
	}
}

func TestSendPushNotification(t *testing.T) {
	service := NewAPNSService()

	// Test with iOS device tokens in recipients
	notification := &struct {
		ID       string
		Type     string
		Content  map[string]interface{}
		Template *struct {
			ID   string
			Data map[string]interface{}
		}
		Recipients  []string
		ScheduledAt *time.Time
	}{
		ID:   "test_notification",
		Type: "test",
		Content: map[string]interface{}{
			"title": "Test Title",
			"body":  "Test Body",
		},
		Recipients: []string{"ios_device_token_123", "ios_device_token_456"},
	}

	// Test sending to iOS devices
	result, err := service.SendPushNotification(context.Background(), notification)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Test with no device tokens
	emptyNotification := &struct {
		ID       string
		Type     string
		Content  map[string]interface{}
		Template *struct {
			ID   string
			Data map[string]interface{}
		}
		Recipients  []string
		ScheduledAt *time.Time
	}{
		ID:   "test_notification",
		Type: "test",
		Content: map[string]interface{}{
			"title": "Test Title",
			"body":  "Test Body",
		},
		Recipients: []string{},
	}

	result, err = service.SendPushNotification(context.Background(), emptyNotification)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Test with invalid notification
	_, err = service.SendPushNotification(context.Background(), "invalid")
	if err != ErrInvalidNotificationPayload {
		t.Errorf("Expected ErrInvalidNotificationPayload, got %v", err)
	}
}
