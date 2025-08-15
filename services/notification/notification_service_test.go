package notification

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
type MockSchedulerService struct{}

func (m *MockEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	return &struct {
		ID      string    `json:"id"`
		Status  string    `json:"status"`
		Message string    `json:"message"`
		SentAt  time.Time `json:"sent_at"`
		Channel string    `json:"channel"`
	}{
		ID:      "test-email-id",
		Status:  "sent",
		Message: "Email sent successfully",
		SentAt:  time.Now(),
		Channel: "email",
	}, nil
}

func (m *MockSlackService) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	return &struct {
		ID      string    `json:"id"`
		Status  string    `json:"status"`
		Message string    `json:"message"`
		SentAt  time.Time `json:"sent_at"`
		Channel string    `json:"channel"`
	}{
		ID:      "test-slack-id",
		Status:  "sent",
		Message: "Slack message sent successfully",
		SentAt:  time.Now(),
		Channel: "slack",
	}, nil
}

func (m *MockInAppService) SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	return &struct {
		ID      string    `json:"id"`
		Status  string    `json:"status"`
		Message string    `json:"message"`
		SentAt  time.Time `json:"sent_at"`
		Channel string    `json:"channel"`
	}{
		ID:      "test-inapp-id",
		Status:  "sent",
		Message: "In-app notification sent successfully",
		SentAt:  time.Now(),
		Channel: "in_app",
	}, nil
}

func (m *MockSchedulerService) ScheduleJob(jobID string, scheduledTime time.Time, job func()) error {
	return nil
}

func (m *MockSchedulerService) CancelJob(jobID string) error {
	return nil
}

func TestNewNotificationManager(t *testing.T) {
	emailService := &MockEmailService{}
	slackService := &MockSlackService{}
	inAppService := &MockInAppService{}
	schedulerService := &MockSchedulerService{}

	manager := NewNotificationManager(emailService, slackService, inAppService, schedulerService)
	assert.NotNil(t, manager)
	assert.Equal(t, emailService, manager.emailService)
	assert.Equal(t, slackService, manager.slackService)
	assert.Equal(t, inAppService, manager.inAppService)
	assert.Equal(t, schedulerService, manager.scheduler)
}

func TestNotificationManager_SendNotification(t *testing.T) {
	emailService := &MockEmailService{}
	slackService := &MockSlackService{}
	inAppService := &MockInAppService{}
	schedulerService := &MockSchedulerService{}

	manager := NewNotificationManager(emailService, slackService, inAppService, schedulerService)

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
				Priority    string
				Title       string
				Message     string
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}{
				ID:         "test-1",
				Type:       "email",
				Priority:   "normal",
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
				Priority    string
				Title       string
				Message     string
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}{
				ID:         "test-2",
				Type:       "slack",
				Priority:   "normal",
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
				Priority    string
				Title       string
				Message     string
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}{
				ID:         "test-3",
				Type:       "in_app",
				Priority:   "normal",
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
	schedulerService := &MockSchedulerService{}

	manager := NewNotificationManager(emailService, slackService, inAppService, schedulerService)
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
