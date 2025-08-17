package validation

import (
	"fmt"
	"testing"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/stretchr/testify/assert"
)

func TestNotificationValidator_ValidateNotificationRequest(t *testing.T) {
	validator := NewNotificationValidator()

	tests := []struct {
		name     string
		request  *models.NotificationRequest
		expected bool
	}{
		{
			name: "Valid email notification",
			request: &models.NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Recipients: []string{"user-123", "user-456"},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: true,
		},
		{
			name: "Valid slack notification",
			request: &models.NotificationRequest{
				Type: "slack",
				Content: map[string]interface{}{
					"text": "Test slack message",
				},
				Recipients: []string{"user-123"},
			},
			expected: true,
		},
		{
			name: "Valid push notification",
			request: &models.NotificationRequest{
				Type: "ios_push",
				Content: map[string]interface{}{
					"title": "Test Title",
					"body":  "Test body",
				},
				Recipients: []string{"user-123"},
			},
			expected: true,
		},
		{
			name: "Valid template notification",
			request: &models.NotificationRequest{
				Type: "email",
				Template: &models.TemplateData{
					ID:      "550e8400-e29b-41d4-a716-446655440000",
					Version: 1,
					Data: map[string]interface{}{
						"name": "John Doe",
					},
				},
				Recipients: []string{"user-123"},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: true,
		},
		{
			name: "Valid scheduled notification",
			request: &models.NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Recipients: []string{"user-123"},
				ScheduledAt: func() *time.Time {
					t := time.Now().Add(time.Hour)
					return &t
				}(),
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: true,
		},
		{
			name: "Invalid - missing type",
			request: &models.NotificationRequest{
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Recipients: []string{"user-123"},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: false,
		},
		{
			name: "Invalid - missing recipients",
			request: &models.NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: false,
		},
		{
			name: "Invalid - email without from field",
			request: &models.NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Recipients: []string{"user-123"},
			},
			expected: false,
		},
		{
			name: "Invalid - non-email with from field",
			request: &models.NotificationRequest{
				Type: "slack",
				Content: map[string]interface{}{
					"text": "Test message",
				},
				Recipients: []string{"user-123"},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: false,
		},
		{
			name: "Invalid - content and template both provided",
			request: &models.NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Template: &models.TemplateData{
					ID:      "550e8400-e29b-41d4-a716-446655440000",
					Version: 1,
					Data: map[string]interface{}{
						"name": "John Doe",
					},
				},
				Recipients: []string{"user-123"},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: false,
		},
		{
			name: "Invalid - neither content nor template provided",
			request: &models.NotificationRequest{
				Type:       "email",
				Recipients: []string{"user-123"},
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: false,
		},
		{
			name: "Invalid - past scheduled time",
			request: &models.NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Recipients: []string{"user-123"},
				ScheduledAt: func() *time.Time {
					t := time.Now().Add(-time.Hour)
					return &t
				}(),
				From: &struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateNotificationRequest(tt.request)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateNotificationRequest() = %v, expected %v", result.IsValid, tt.expected)
				if !result.IsValid {
					t.Logf("Validation errors: %+v", result.Errors)
				}
			}
		})
	}
}

func TestNotificationValidator_ValidateRecipients(t *testing.T) {
	validator := NewNotificationValidator()

	tests := []struct {
		name       string
		recipients []string
		expected   bool
	}{
		{
			name:       "Valid recipients",
			recipients: []string{"user-123", "user-456"},
			expected:   true,
		},
		{
			name:       "Empty recipients",
			recipients: []string{},
			expected:   false,
		},
		{
			name:       "Empty recipient string",
			recipients: []string{""},
			expected:   false,
		},
		{
			name:       "Whitespace only recipient",
			recipients: []string{"   "},
			expected:   false,
		},
		{
			name:       "Invalid characters in recipient",
			recipients: []string{"user@123"},
			expected:   false,
		},
		{
			name: "Too many recipients",
			recipients: func() []string {
				recipients := make([]string, 1001)
				for i := range recipients {
					recipients[i] = fmt.Sprintf("user-%d", i)
				}
				return recipients
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.validateRecipients(tt.recipients)
			isValid := len(errors) == 0
			if isValid != tt.expected {
				t.Errorf("validateRecipients() = %v, expected %v", isValid, tt.expected)
				if !isValid {
					t.Logf("Validation errors: %+v", errors)
				}
			}
		})
	}
}

func TestNotificationValidator_ValidateTemplateVersion(t *testing.T) {
	validator := NewNotificationValidator()

	tests := []struct {
		name     string
		template *models.TemplateData
		expected bool
	}{
		{
			name: "Valid template with version",
			template: &models.TemplateData{
				ID:      "550e8400-e29b-41d4-a716-446655440000",
				Version: 1,
				Data: map[string]interface{}{
					"name": "John Doe",
				},
			},
			expected: true,
		},
		{
			name: "Invalid template without version (version is now mandatory)",
			template: &models.TemplateData{
				ID: "550e8400-e29b-41d4-a716-446655440000",
				Data: map[string]interface{}{
					"name": "John Doe",
				},
			},
			expected: false,
		},
		{
			name: "Invalid template with negative version",
			template: &models.TemplateData{
				ID:      "550e8400-e29b-41d4-a716-446655440000",
				Version: -1,
				Data: map[string]interface{}{
					"name": "John Doe",
				},
			},
			expected: false,
		},
		{
			name: "Invalid template with zero version (version must be positive)",
			template: &models.TemplateData{
				ID:      "550e8400-e29b-41d4-a716-446655440000",
				Version: 0,
				Data: map[string]interface{}{
					"name": "John Doe",
				},
			},
			expected: false, // Zero version is now invalid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.validateTemplate(tt.template)
			isValid := len(errors) == 0
			if isValid != tt.expected {
				t.Errorf("validateTemplate() = %v, expected %v", isValid, tt.expected)
				if !isValid {
					t.Logf("Validation errors: %+v", errors)
				}
			}
		})
	}
}

func TestNotificationValidator_ValidateNotificationID(t *testing.T) {
	validator := NewNotificationValidator()

	tests := []struct {
		name           string
		notificationID string
		expectedValid  bool
		expectedErrors []string
	}{
		{
			name:           "valid UUID",
			notificationID: "123e4567-e89b-12d3-a456-426614174000",
			expectedValid:  true,
			expectedErrors: []string{},
		},
		{
			name:           "valid UUID uppercase",
			notificationID: "123E4567-E89B-12D3-A456-426614174000",
			expectedValid:  true,
			expectedErrors: []string{},
		},
		{
			name:           "empty notification ID",
			notificationID: "",
			expectedValid:  false,
			expectedErrors: []string{"notification ID is required"},
		},
		{
			name:           "invalid UUID format",
			notificationID: "invalid-uuid-format",
			expectedValid:  false,
			expectedErrors: []string{"notification ID must be a valid UUID"},
		},
		{
			name:           "too long notification ID",
			notificationID: "123e4567-e89b-12d3-a456-426614174000-extra",
			expectedValid:  false,
			expectedErrors: []string{"notification ID must be a valid UUID", "notification ID cannot exceed 36 characters"},
		},
		{
			name:           "partial UUID",
			notificationID: "123e4567-e89b-12d3",
			expectedValid:  false,
			expectedErrors: []string{"notification ID must be a valid UUID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateNotificationID(tt.notificationID)

			assert.Equal(t, tt.expectedValid, result.IsValid, "Expected valid to be %v, got %v", tt.expectedValid, result.IsValid)

			if len(tt.expectedErrors) > 0 {
				assert.Len(t, result.Errors, len(tt.expectedErrors), "Expected %d errors, got %d", len(tt.expectedErrors), len(result.Errors))

				for i, expectedError := range tt.expectedErrors {
					assert.Contains(t, result.Errors[i].Message, expectedError, "Expected error message to contain '%s'", expectedError)
				}
			} else {
				assert.Empty(t, result.Errors, "Expected no errors, got %d", len(result.Errors))
			}
		})
	}
}
