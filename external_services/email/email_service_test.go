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
		errorMessage string
	}{
		{
			name: "valid email notification",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: false,
		},
		{
			name: "missing ID",
			notification: &models.EmailNotificationRequest{
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email notification ID is required",
		},
		{
			name: "missing type",
			notification: &models.EmailNotificationRequest{
				ID: "test-123",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email notification type is required",
		},
		{
			name: "missing content",
			notification: &models.EmailNotificationRequest{
				ID:        "test-123",
				Type:      "email",
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email subject is required",
		},
		{
			name: "missing recipients",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "recipient is required",
		},
		{
			name: "empty recipients list",
			notification: &models.EmailNotificationRequest{
				ID:        "test-123",
				Type:      "email",
				Content:   models.EmailContent{},
				Recipient: "",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email subject is required",
		},
		{
			name: "missing from email",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
			},
			expectError: false, // From is optional
		},
		{
			name: "invalid from email",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "invalid-email",
				},
			},
			expectError:  true,
			errorMessage: "invalid from email address: invalid-email",
		},
		{
			name: "invalid recipient email",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "Test Body",
				},
				Recipient: "invalid-email",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError: true,
		},
		{
			name: "missing subject",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email subject is required",
		},
		{
			name: "missing body",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject: "Test Subject",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email body is required",
		},
		{
			name: "empty subject",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "",
					EmailBody: "Test Body",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email subject is required",
		},
		{
			name: "empty body",
			notification: &models.EmailNotificationRequest{
				ID:   "test-123",
				Type: "email",
				Content: models.EmailContent{
					Subject:   "Test Subject",
					EmailBody: "",
				},
				Recipient: "test@example.com",
				From: &models.EmailSender{
					Email: "sender@example.com",
				},
			},
			expectError:  true,
			errorMessage: "email body is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := models.ValidateEmailNotification(tt.notification)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewEmailNotification(t *testing.T) {
	t.Run("valid notification", func(t *testing.T) {
		content := models.EmailContent{
			Subject:   "Test Subject",
			EmailBody: "Test Body",
		}
		recipient := "test@example.com"

		notification := &models.EmailNotificationRequest{
			ID:        "test-123",
			Type:      "email",
			Content:   content,
			Recipient: recipient,
			From: &models.EmailSender{
				Email: "sender@example.com",
			},
		}

		err := models.ValidateEmailNotification(notification)
		require.NoError(t, err)
		assert.Equal(t, "test-123", notification.ID)
		assert.Equal(t, "email", notification.Type)
		assert.Equal(t, content, notification.Content)
		assert.Equal(t, recipient, notification.Recipient)
		assert.Equal(t, "sender@example.com", notification.From.Email)
	})

	t.Run("invalid notification", func(t *testing.T) {
		content := models.EmailContent{
			Subject: "Test Subject",
			// missing email_body
		}
		recipient := "test@example.com"

		notification := &models.EmailNotificationRequest{
			ID:        "test-123",
			Type:      "email",
			Content:   content,
			Recipient: recipient,
			From: &models.EmailSender{
				Email: "sender@example.com",
			},
		}

		err := models.ValidateEmailNotification(notification)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email body is required")
	})
}

func TestNewEmailNotificationWithTemplate(t *testing.T) {
	t.Run("valid notification with template", func(t *testing.T) {
		content := models.EmailContent{
			Subject:   "Test Subject",
			EmailBody: "Test Body",
		}
		recipient := "test@example.com"

		notification := &models.EmailNotificationRequest{
			ID:        "test-123",
			Type:      "email",
			Content:   content,
			Recipient: recipient,
			From: &models.EmailSender{
				Email: "sender@example.com",
			},
		}

		err := models.ValidateEmailNotification(notification)
		require.NoError(t, err)
		assert.Equal(t, "test-123", notification.ID)
		assert.Equal(t, "email", notification.Type)
		assert.Equal(t, content, notification.Content)
		assert.Equal(t, recipient, notification.Recipient)
		assert.Equal(t, "sender@example.com", notification.From.Email)
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
			err := models.ValidateEmailAddress(tt.email)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
