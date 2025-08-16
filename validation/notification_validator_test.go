package validation

import (
	"fmt"
	"testing"
	"time"

	"github.com/gaurav2721/notification-service/models"
)

func TestNotificationValidator_ValidateNotificationRequest(t *testing.T) {
	validator := NewNotificationValidator()

	tests := []struct {
		name     string
		request  *NotificationRequest
		expected bool
	}{
		{
			name: "Valid email notification",
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
				Type: "email",
				Template: &models.TemplateData{
					ID: "550e8400-e29b-41d4-a716-446655440000",
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test email body",
				},
				Template: &models.TemplateData{
					ID: "550e8400-e29b-41d4-a716-446655440000",
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
			request: &NotificationRequest{
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
			request: &NotificationRequest{
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
