package validation

import (
	"testing"

	"github.com/gaurav2721/notification-service/models"
)

func TestTemplateValidator_ValidateTemplateRequest(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name     string
		request  *models.TemplateRequest
		expected bool
	}{
		{
			name: "valid email template",
			request: &models.TemplateRequest{
				Name: "Welcome Email",
				Type: models.EmailNotification,
				Content: models.TemplateContent{
					Subject:   "Welcome to our service",
					EmailBody: "Hello {{name}}, welcome to our platform!",
				},
				RequiredVariables: []string{"name"},
				Description:       "Welcome email template",
			},
			expected: true,
		},
		{
			name: "valid slack template",
			request: &models.TemplateRequest{
				Name: "Alert Notification",
				Type: models.SlackNotification,
				Content: models.TemplateContent{
					Text: "Alert: {{message}} for {{user}}",
				},
				RequiredVariables: []string{"message", "user"},
				Description:       "Slack alert template",
			},
			expected: true,
		},
		{
			name: "valid in-app template",
			request: &models.TemplateRequest{
				Name: "Push Notification",
				Type: models.InAppNotification,
				Content: models.TemplateContent{
					Title: "New Message",
					Body:  "You have a new message from {{sender}}",
				},
				RequiredVariables: []string{"sender"},
				Description:       "In-app notification template",
			},
			expected: true,
		},
		{
			name: "missing name",
			request: &models.TemplateRequest{
				Type: models.EmailNotification,
				Content: models.TemplateContent{
					Subject:   "Welcome",
					EmailBody: "Hello",
				},
				RequiredVariables: []string{"name"},
			},
			expected: false,
		},
		{
			name: "invalid template type",
			request: &models.TemplateRequest{
				Name: "Test Template",
				Type: "invalid_type",
				Content: models.TemplateContent{
					Subject:   "Welcome",
					EmailBody: "Hello",
				},
				RequiredVariables: []string{"name"},
			},
			expected: false,
		},
		{
			name: "missing email subject",
			request: &models.TemplateRequest{
				Name: "Email Template",
				Type: models.EmailNotification,
				Content: models.TemplateContent{
					EmailBody: "Hello {{name}}",
				},
				RequiredVariables: []string{"name"},
			},
			expected: false,
		},
		{
			name: "missing slack text",
			request: &models.TemplateRequest{
				Name:              "Slack Template",
				Type:              models.SlackNotification,
				Content:           models.TemplateContent{},
				RequiredVariables: []string{"message"},
			},
			expected: false,
		},
		{
			name: "missing in-app title",
			request: &models.TemplateRequest{
				Name: "In-App Template",
				Type: models.InAppNotification,
				Content: models.TemplateContent{
					Body: "Hello {{name}}",
				},
				RequiredVariables: []string{"name"},
			},
			expected: false,
		},
		{
			name: "empty required variables",
			request: &models.TemplateRequest{
				Name: "Template",
				Type: models.EmailNotification,
				Content: models.TemplateContent{
					Subject:   "Welcome",
					EmailBody: "Hello",
				},
				RequiredVariables: []string{},
			},
			expected: false,
		},
		{
			name: "duplicate variables",
			request: &models.TemplateRequest{
				Name: "Template",
				Type: models.EmailNotification,
				Content: models.TemplateContent{
					Subject:   "Welcome",
					EmailBody: "Hello",
				},
				RequiredVariables: []string{"name", "name"},
			},
			expected: false,
		},
		{
			name: "invalid variable name",
			request: &models.TemplateRequest{
				Name: "Template",
				Type: models.EmailNotification,
				Content: models.TemplateContent{
					Subject:   "Welcome",
					EmailBody: "Hello",
				},
				RequiredVariables: []string{"123invalid"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateTemplateRequest(tt.request)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateTemplateRequest() = %v, expected %v", result.IsValid, tt.expected)
				if !result.IsValid {
					t.Logf("Validation errors: %+v", result.Errors)
				}
			}
		})
	}
}

func TestTemplateValidator_ValidateTemplateData(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name     string
		data     *models.TemplateData
		expected bool
	}{
		{
			name: "valid template data",
			data: &models.TemplateData{
				ID:      "550e8400-e29b-41d4-a716-446655440000",
				Version: 1,
				Data: map[string]interface{}{
					"name": "John",
				},
			},
			expected: true,
		},
		{
			name:     "nil template data",
			data:     nil,
			expected: false,
		},
		{
			name: "empty template ID",
			data: &models.TemplateData{
				ID:      "",
				Version: 1,
				Data:    map[string]interface{}{},
			},
			expected: false,
		},
		{
			name: "invalid version",
			data: &models.TemplateData{
				ID:      "550e8400-e29b-41d4-a716-446655440000",
				Version: 0,
				Data:    map[string]interface{}{},
			},
			expected: false,
		},
		{
			name: "nil data",
			data: &models.TemplateData{
				ID:      "550e8400-e29b-41d4-a716-446655440000",
				Version: 1,
				Data:    nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateTemplateData(tt.data)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateTemplateData() = %v, expected %v", result.IsValid, tt.expected)
				if !result.IsValid {
					t.Logf("Validation errors: %+v", result.Errors)
				}
			}
		})
	}
}

func TestTemplateValidator_ValidateTemplateID(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name       string
		templateID string
		expected   bool
	}{
		{
			name:       "valid UUID",
			templateID: "550e8400-e29b-41d4-a716-446655440000",
			expected:   true,
		},
		{
			name:       "empty template ID",
			templateID: "",
			expected:   false,
		},
		{
			name:       "invalid UUID format",
			templateID: "invalid-uuid",
			expected:   false,
		},
		{
			name:       "UUID with uppercase",
			templateID: "550E8400-E29B-41D4-A716-446655440000",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateTemplateID(tt.templateID)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateTemplateID() = %v, expected %v", result.IsValid, tt.expected)
				if !result.IsValid {
					t.Logf("Validation errors: %+v", result.Errors)
				}
			}
		})
	}
}

func TestTemplateValidator_ValidateTemplateVersion(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name     string
		version  int
		expected bool
	}{
		{
			name:     "valid version",
			version:  1,
			expected: true,
		},
		{
			name:     "valid version greater than 1",
			version:  5,
			expected: true,
		},
		{
			name:     "zero version",
			version:  0,
			expected: false,
		},
		{
			name:     "negative version",
			version:  -1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateTemplateVersion(tt.version)
			if result.IsValid != tt.expected {
				t.Errorf("ValidateTemplateVersion() = %v, expected %v", result.IsValid, tt.expected)
				if !result.IsValid {
					t.Logf("Validation errors: %+v", result.Errors)
				}
			}
		})
	}
}

func TestTemplateValidator_validateTemplateName(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name         string
		templateName string
		expected     bool
	}{
		{
			name:         "valid name",
			templateName: "Welcome Email Template",
			expected:     true,
		},
		{
			name:         "valid name with hyphens",
			templateName: "welcome-email-template",
			expected:     true,
		},
		{
			name:         "valid name with underscores",
			templateName: "welcome_email_template",
			expected:     true,
		},
		{
			name:         "empty name",
			templateName: "",
			expected:     false,
		},
		{
			name:         "whitespace only",
			templateName: "   ",
			expected:     false,
		},
		{
			name:         "name too long",
			templateName: "This is a very long template name that exceeds the maximum allowed length of one hundred characters and should fail validation",
			expected:     false,
		},
		{
			name:         "invalid characters",
			templateName: "Template@#$%",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.validateTemplateName(tt.templateName)
			isValid := len(errors) == 0
			if isValid != tt.expected {
				t.Errorf("validateTemplateName() = %v, expected %v", isValid, tt.expected)
				if !isValid {
					t.Logf("Validation errors: %+v", errors)
				}
			}
		})
	}
}

func TestTemplateValidator_validateRequiredVariables(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name      string
		variables []string
		expected  bool
	}{
		{
			name:      "valid variables",
			variables: []string{"name", "email", "company"},
			expected:  true,
		},
		{
			name:      "nil variables",
			variables: nil,
			expected:  false,
		},
		{
			name:      "empty variables",
			variables: []string{},
			expected:  false,
		},
		{
			name:      "duplicate variables",
			variables: []string{"name", "email", "name"},
			expected:  false,
		},
		{
			name:      "empty variable name",
			variables: []string{"name", "", "email"},
			expected:  false,
		},
		{
			name:      "invalid variable name starting with number",
			variables: []string{"name", "1email"},
			expected:  false,
		},
		{
			name:      "invalid variable name with special characters",
			variables: []string{"name", "email@domain"},
			expected:  false,
		},
		{
			name:      "valid variable starting with underscore",
			variables: []string{"_name", "email"},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.validateRequiredVariables(tt.variables)
			isValid := len(errors) == 0
			if isValid != tt.expected {
				t.Errorf("validateRequiredVariables() = %v, expected %v", isValid, tt.expected)
				if !isValid {
					t.Logf("Validation errors: %+v", errors)
				}
			}
		})
	}
}
