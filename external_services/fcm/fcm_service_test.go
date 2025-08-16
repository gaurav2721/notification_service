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
		ProjectID: "test_project",
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

func TestRegisterDeviceToken(t *testing.T) {
	service := NewFCMService()

	// Test valid registration
	err := service.RegisterDeviceToken("user1", "device_token_1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test duplicate registration
	err = service.RegisterDeviceToken("user1", "device_token_1")
	if err != nil {
		t.Errorf("Expected no error for duplicate registration, got %v", err)
	}

	// Test invalid parameters
	err = service.RegisterDeviceToken("", "device_token_1")
	if err != ErrInvalidDeviceToken {
		t.Errorf("Expected ErrInvalidDeviceToken, got %v", err)
	}

	err = service.RegisterDeviceToken("user1", "")
	if err != ErrInvalidDeviceToken {
		t.Errorf("Expected ErrInvalidDeviceToken, got %v", err)
	}
}

func TestGetDeviceTokensForUser(t *testing.T) {
	service := NewFCMService()

	// Register tokens
	service.RegisterDeviceToken("user1", "device_token_1")
	service.RegisterDeviceToken("user1", "device_token_2")
	service.RegisterDeviceToken("user2", "device_token_3")

	// Test getting tokens for user1
	tokens, err := service.GetDeviceTokensForUser("user1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tokens) != 2 {
		t.Errorf("Expected 2 tokens, got %d", len(tokens))
	}

	// Test getting tokens for non-existent user
	tokens, err = service.GetDeviceTokensForUser("user3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tokens) != 0 {
		t.Errorf("Expected 0 tokens, got %d", len(tokens))
	}

	// Test invalid user ID
	_, err = service.GetDeviceTokensForUser("")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUnregisterDeviceToken(t *testing.T) {
	service := NewFCMService()

	// Register token
	service.RegisterDeviceToken("user1", "device_token_1")

	// Test successful unregistration
	err := service.UnregisterDeviceToken("user1", "device_token_1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify token is removed
	tokens, _ := service.GetDeviceTokensForUser("user1")
	if len(tokens) != 0 {
		t.Errorf("Expected 0 tokens after unregistration, got %d", len(tokens))
	}

	// Test unregistering non-existent token
	err = service.UnregisterDeviceToken("user1", "device_token_1")
	if err != ErrDeviceTokenNotFound {
		t.Errorf("Expected ErrDeviceTokenNotFound, got %v", err)
	}

	// Test unregistering from non-existent user
	err = service.UnregisterDeviceToken("user2", "device_token_1")
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestSendPushNotification(t *testing.T) {
	service := NewFCMService()

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
		Recipients: []string{"user1"},
	}

	// Test sending to user with no device tokens
	result, err := service.SendPushNotification(context.Background(), notification)
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
