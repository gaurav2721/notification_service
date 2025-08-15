package notification_manager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock services for testing
type MockEmailService struct{}
type MockSlackService struct{}
type MockInAppService struct{}

func (m *MockEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	return &struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     "email-123",
		Status: "sent",
	}, nil
}

func (m *MockSlackService) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	return &struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     "slack-123",
		Status: "sent",
	}, nil
}

func (m *MockInAppService) SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	return &struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     "inapp-123",
		Status: "sent",
	}, nil
}

func TestNewNotificationManager(t *testing.T) {
	emailService := &MockEmailService{}
	slackService := &MockSlackService{}
	inAppService := &MockInAppService{}

	manager := NewNotificationManager(emailService, slackService, inAppService)
	assert.NotNil(t, manager)
}

func TestNotificationManager_SendNotification(t *testing.T) {
	emailService := &MockEmailService{}
	slackService := &MockSlackService{}
	inAppService := &MockInAppService{}

	manager := NewNotificationManager(emailService, slackService, inAppService)

	tests := []struct {
		name         string
		notification interface{}
		expectError  bool
	}{
		{
			name: "Send Email Notification",
			notification: &struct {
				ID          string
				Type        string
				Title       string
				Message     string
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}{
				ID:         "test-1",
				Type:       "email",
				Title:      "Test Email",
				Message:    "This is a test email",
				Recipients: []string{"test@example.com"},
			},
			expectError: false,
		},
		{
			name: "Send Slack Notification",
			notification: &struct {
				ID          string
				Type        string
				Title       string
				Message     string
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}{
				ID:         "test-2",
				Type:       "slack",
				Title:      "Test Slack",
				Message:    "This is a test slack message",
				Recipients: []string{"#general"},
			},
			expectError: false,
		},
		{
			name: "Send In-App Notification",
			notification: &struct {
				ID          string
				Type        string
				Title       string
				Message     string
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}{
				ID:         "test-3",
				Type:       "in_app",
				Title:      "Test In-App",
				Message:    "This is a test in-app notification",
				Recipients: []string{"user-1"},
			},
			expectError: false,
		},
		{
			name:         "Unsupported Notification Type",
			notification: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := manager.SendNotification(ctx, tt.notification)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestNotificationManager_TemplateOperations(t *testing.T) {
	emailService := &MockEmailService{}
	slackService := &MockSlackService{}
	inAppService := &MockInAppService{}

	manager := NewNotificationManager(emailService, slackService, inAppService)
	ctx := context.Background()

	t.Run("Create Template", func(t *testing.T) {
		template := &struct {
			ID   string
			Name string
			Type string
			Body string
		}{
			ID:   "template-1",
			Name: "Welcome Email",
			Type: "email",
			Body: "Welcome to our platform!",
		}

		err := manager.CreateTemplate(ctx, template)
		assert.NoError(t, err)
	})

	t.Run("Get Template", func(t *testing.T) {
		template, err := manager.GetTemplate(ctx, "template-1")
		assert.NoError(t, err)
		assert.NotNil(t, template)
	})

	t.Run("Update Template", func(t *testing.T) {
		template := &struct {
			ID   string
			Name string
			Type string
			Body string
		}{
			ID:   "template-1",
			Name: "Updated Welcome Email",
			Type: "email",
			Body: "Updated welcome message!",
		}

		err := manager.UpdateTemplate(ctx, template)
		assert.NoError(t, err)
	})

	t.Run("Delete Template", func(t *testing.T) {
		err := manager.DeleteTemplate(ctx, "template-1")
		assert.NoError(t, err)

		// Verify template is deleted
		_, err = manager.GetTemplate(ctx, "template-1")
		assert.Error(t, err)
	})
}
