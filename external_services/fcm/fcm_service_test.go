package fcm

import (
	"context"
	"testing"
	"time"
)

func TestNewFCMService(t *testing.T) {
	service := NewFCMService()
	if service == nil {
		t.Fatal("Expected FCM service to be created, got nil")
	}
}

func TestNewFCMServiceWithConfig(t *testing.T) {
	// Test with valid config
	config := &FCMConfig{
		ServerKey: "test_server_key",
		Timeout:   30,
		BatchSize: 1000,
	}

	service, err := NewFCMServiceWithConfig(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if service == nil {
		t.Fatal("Expected FCM service to be created, got nil")
	}

	// Test with nil config
	_, err = NewFCMServiceWithConfig(nil)
	if err != ErrInvalidConfiguration {
		t.Errorf("Expected ErrInvalidConfiguration, got %v", err)
	}

	// Test with empty server key
	config.ServerKey = ""
	_, err = NewFCMServiceWithConfig(config)
	if err != ErrInvalidConfiguration {
		t.Errorf("Expected ErrInvalidConfiguration, got %v", err)
	}
}

func TestSendPushNotification(t *testing.T) {
	service := NewFCMService()

	// Test with Android device tokens in recipients
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
			"data": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		Recipients: []string{"android_device_token_123", "android_device_token_456"},
	}

	// Test sending to Android devices
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
