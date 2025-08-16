package consumers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailProcessor_ProcessNotification(t *testing.T) {
	// Create email processor
	processor := NewEmailProcessor()
	require.NotNil(t, processor)

	// Create test notification data
	notificationData := map[string]interface{}{
		"notification_id": "test-123",
		"type":            "email",
		"content": map[string]interface{}{
			"subject":    "Test Email Subject",
			"email_body": "<h1>Test Email Body</h1><p>This is a test email.</p>",
		},
		"recipient": map[string]interface{}{
			"user_id":   "user-123",
			"email":     "test@example.com",
			"full_name": "Test User",
		},
		"from": map[string]interface{}{
			"email": "sender@example.com",
		},
		"created_at": time.Now(),
	}

	// Convert to JSON
	payload, err := json.Marshal(notificationData)
	require.NoError(t, err)

	// Create notification message
	message := NotificationMessage{
		Type:      EmailNotification,
		Payload:   string(payload),
		ID:        "test-123",
		Timestamp: time.Now().Unix(),
	}

	// Process notification
	ctx := context.Background()
	err = processor.ProcessNotification(ctx, message)

	// Since we're using a mock email service in test environment (no SMTP config),
	// the email should be processed successfully
	assert.NoError(t, err)
}

func TestEmailProcessor_ProcessNotification_InvalidPayload(t *testing.T) {
	// Create email processor
	processor := NewEmailProcessor()
	require.NotNil(t, processor)

	// Create invalid notification message
	message := NotificationMessage{
		Type:      EmailNotification,
		Payload:   "invalid json",
		ID:        "test-123",
		Timestamp: time.Now().Unix(),
	}

	// Process notification
	ctx := context.Background()
	err := processor.ProcessNotification(ctx, message)

	// Should fail due to invalid JSON
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse notification payload")
}

func TestEmailProcessor_ProcessNotification_MissingEmail(t *testing.T) {
	// Create email processor
	processor := NewEmailProcessor()
	require.NotNil(t, processor)

	// Create notification data without email
	notificationData := map[string]interface{}{
		"notification_id": "test-123",
		"type":            "email",
		"content": map[string]interface{}{
			"subject":    "Test Email Subject",
			"email_body": "<h1>Test Email Body</h1>",
		},
		"recipient": map[string]interface{}{
			"user_id":   "user-123",
			"full_name": "Test User",
			// Missing email field
		},
		"created_at": time.Now(),
	}

	// Convert to JSON
	payload, err := json.Marshal(notificationData)
	require.NoError(t, err)

	// Create notification message
	message := NotificationMessage{
		Type:      EmailNotification,
		Payload:   string(payload),
		ID:        "test-123",
		Timestamp: time.Now().Unix(),
	}

	// Process notification
	ctx := context.Background()
	err = processor.ProcessNotification(ctx, message)

	// Should fail due to missing email
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing or invalid email address")
}

func TestEmailProcessor_GetNotificationType(t *testing.T) {
	// Create email processor
	processor := NewEmailProcessor()
	require.NotNil(t, processor)

	// Check notification type
	notificationType := processor.GetNotificationType()
	assert.Equal(t, EmailNotification, notificationType)
}

func TestEmailProcessor_WithCustomEmailService(t *testing.T) {
	// Create a mock email service for testing
	mockEmailService := &mockEmailService{}

	// Create email processor with custom service
	processor := NewEmailProcessorWithService(mockEmailService)
	require.NotNil(t, processor)

	// Create test notification data
	notificationData := map[string]interface{}{
		"notification_id": "test-123",
		"type":            "email",
		"content": map[string]interface{}{
			"subject":    "Test Email Subject",
			"email_body": "<h1>Test Email Body</h1>",
		},
		"recipient": map[string]interface{}{
			"user_id": "user-123",
			"email":   "test@example.com",
		},
		"created_at": time.Now(),
	}

	// Convert to JSON
	payload, err := json.Marshal(notificationData)
	require.NoError(t, err)

	// Create notification message
	message := NotificationMessage{
		Type:      EmailNotification,
		Payload:   string(payload),
		ID:        "test-123",
		Timestamp: time.Now().Unix(),
	}

	// Process notification
	ctx := context.Background()
	err = processor.ProcessNotification(ctx, message)

	// Should succeed with mock service
	assert.NoError(t, err)
	assert.True(t, mockEmailService.sendEmailCalled)
}

// mockEmailService is a mock implementation for testing
type mockEmailService struct {
	sendEmailCalled bool
}

func (m *mockEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	m.sendEmailCalled = true

	// Return a mock response
	return &models.EmailResponse{
		ID:      "test-123",
		Status:  "sent",
		Message: "Email sent successfully",
		SentAt:  time.Now(),
		Channel: "email",
	}, nil
}
