package email

import (
	"testing"

	"github.com/gaurav2721/notification-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmailNotification(t *testing.T) {
	tests := []struct {
		name         string
		notification *models.EmailNotificationRequest
		expectError  bool
		errorType    error
	}{
		{
			name: "valid email notification",
			notification: &models.EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: false,
		},
		{
			name: "missing ID",
			notification: &EmailNotification{
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingID,
		},
		{
			name: "missing type",
			notification: &EmailNotification{
				ID: "test-123",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingType,
		},
		{
			name: "missing content",
			notification: &EmailNotification{
				ID:         "test-123",
				Type:       "email",
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingContent,
		},
		{
			name: "missing recipients",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrEmptyRecipients,
		},
		{
			name: "empty recipients list",
			notification: &EmailNotification{
				ID:         "test-123",
				Type:       "email",
				Content:    map[string]interface{}{},
				Recipients: []string{},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrEmptyRecipients,
		},
		{
			name: "missing from email",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
			},
			expectError: true,
			errorType:   ErrMissingFromEmail,
		},
		{
			name: "invalid from email",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "invalid-email",
				},
			},
			expectError: true,
			errorType:   ErrInvalidEmail,
		},
		{
			name: "invalid recipient email",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "Test Body",
				},
				Recipients: []string{"invalid-email"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
		},
		{
			name: "missing subject",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingSubject,
		},
		{
			name: "missing body",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject": "Test Subject",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingBody,
		},
		{
			name: "empty subject",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "",
					"email_body": "Test Body",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingSubject,
		},
		{
			name: "empty body",
			notification: &EmailNotification{
				ID:   "test-123",
				Type: "email",
				Content: map[string]interface{}{
					"subject":    "Test Subject",
					"email_body": "",
				},
				Recipients: []string{"test@example.com"},
				From: &EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
			errorType:   ErrMissingBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailNotification(tt.notification)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewEmailNotification(t *testing.T) {
	t.Run("valid notification", func(t *testing.T) {
		content := map[string]interface{}{
			"subject":    "Test Subject",
			"email_body": "Test Body",
		}
		recipients := []string{"test@example.com"}

		notification, err := NewEmailNotification("test-123", "email", content, recipients, "sender@example.com")

		require.NoError(t, err)
		assert.Equal(t, "test-123", notification.ID)
		assert.Equal(t, "email", notification.Type)
		assert.Equal(t, content, notification.Content)
		assert.Equal(t, recipients, notification.Recipients)
		assert.Equal(t, "sender@example.com", notification.From.Email)
	})

	t.Run("invalid notification", func(t *testing.T) {
		content := map[string]interface{}{
			"subject": "Test Subject",
			// missing email_body
		}
		recipients := []string{"test@example.com"}

		notification, err := NewEmailNotification("test-123", "email", content, recipients, "sender@example.com")

		assert.Error(t, err)
		assert.Nil(t, notification)
		assert.ErrorIs(t, err, ErrMissingBody)
	})
}

func TestNewEmailNotificationWithTemplate(t *testing.T) {
	t.Run("valid notification with template", func(t *testing.T) {
		content := map[string]interface{}{
			"subject":    "Test Subject",
			"email_body": "Test Body",
		}
		recipients := []string{"test@example.com"}
		template := &EmailTemplate{
			ID: "welcome-template",
			Data: map[string]interface{}{
				"user_name": "John Doe",
			},
		}

		notification, err := NewEmailNotificationWithTemplate("test-123", "email", content, recipients, "sender@example.com", template)

		require.NoError(t, err)
		assert.Equal(t, "test-123", notification.ID)
		assert.Equal(t, "email", notification.Type)
		assert.Equal(t, content, notification.Content)
		assert.Equal(t, recipients, notification.Recipients)
		assert.Equal(t, "sender@example.com", notification.From.Email)
		assert.Equal(t, template, notification.Template)
	})
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@sub.example.com", true},
		{"valid email with plus", "test+tag@example.com", true},
		{"valid email with dots", "test.name@example.com", true},
		{"invalid email - no @", "testexample.com", false},
		{"invalid email - no domain", "test@", false},
		{"invalid email - no local part", "@example.com", false},
		{"invalid email - spaces", "test @example.com", false},
		{"invalid email - empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}
