package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/gaurav2721/notification-service/models"
)

// NotificationValidator provides validation methods for notification requests
type NotificationValidator struct{}

// NewNotificationValidator creates a new notification validator
func NewNotificationValidator() *NotificationValidator {
	return &NotificationValidator{}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// ValidateNotificationRequest validates a complete notification request
func (v *NotificationValidator) ValidateNotificationRequest(request *models.NotificationRequest) ValidationResult {
	var errors []ValidationError

	// Validate basic required fields
	if errors = append(errors, v.validateType(request.Type)...); len(errors) > 0 {
		return ValidationResult{IsValid: false, Errors: errors}
	}

	if errors = append(errors, v.validateRecipients(request.Recipients)...); len(errors) > 0 {
		return ValidationResult{IsValid: false, Errors: errors}
	}

	// Validate content vs template (mutual exclusivity)
	if contentErrors := v.validateContentAndTemplate(request.Content, request.Template); len(contentErrors) > 0 {
		errors = append(errors, contentErrors...)
	}

	// Validate content based on type
	if request.Content != nil {
		if contentErrors := v.validateContentByType(request.Type, request.Content); len(contentErrors) > 0 {
			errors = append(errors, contentErrors...)
		}
	}

	// Validate template if provided
	if request.Template != nil {
		if templateErrors := v.validateTemplate(request.Template); len(templateErrors) > 0 {
			errors = append(errors, templateErrors...)
		}
	}

	// Validate from field based on type
	if fromErrors := v.validateFromField(request.Type, request.From); len(fromErrors) > 0 {
		errors = append(errors, fromErrors...)
	}

	// Validate scheduled_at if provided
	if request.ScheduledAt != nil {
		if scheduleErrors := v.validateScheduledAt(*request.ScheduledAt); len(scheduleErrors) > 0 {
			errors = append(errors, scheduleErrors...)
		}
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// validateType validates the notification type
func (v *NotificationValidator) validateType(notificationType string) []ValidationError {
	var errors []ValidationError

	if notificationType == "" {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: "notification type is required",
		})
		return errors
	}

	validTypes := map[string]bool{
		"email":        true,
		"slack":        true,
		"ios_push":     true,
		"android_push": true,
		"in_app":       true,
	}

	if !validTypes[notificationType] {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("invalid notification type: %s. Valid types are: email, slack, ios_push, android_push, in_app", notificationType),
		})
	}

	return errors
}

// validateRecipients validates the recipients array
func (v *NotificationValidator) validateRecipients(recipients []string) []ValidationError {
	var errors []ValidationError

	if len(recipients) == 0 {
		errors = append(errors, ValidationError{
			Field:   "recipients",
			Message: "at least one recipient is required",
		})
		return errors
	}

	// Check for maximum recipients limit
	if len(recipients) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "recipients",
			Message: "maximum 1000 recipients allowed per notification",
		})
	}

	// Validate each recipient
	for i, recipient := range recipients {
		if recipient == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: "recipient cannot be empty",
			})
			continue
		}

		// Trim whitespace
		recipient = strings.TrimSpace(recipient)
		if recipient == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: "recipient cannot be empty after trimming whitespace",
			})
			continue
		}

		// Check for minimum length
		if len(recipient) < 1 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: "recipient must be at least 1 character long",
			})
			continue
		}

		// Check for maximum length
		if len(recipient) > 255 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: "recipient cannot exceed 255 characters",
			})
			continue
		}

		// Check for valid characters (alphanumeric, hyphens, underscores)
		validRecipientRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !validRecipientRegex.MatchString(recipient) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: "recipient can only contain alphanumeric characters, hyphens, and underscores",
			})
		}
	}

	return errors
}

// validateContentAndTemplate validates that either content or template is provided, but not both
func (v *NotificationValidator) validateContentAndTemplate(content map[string]interface{}, template *models.TemplateData) []ValidationError {
	var errors []ValidationError

	hasContent := content != nil && len(content) > 0
	hasTemplate := template != nil

	if !hasContent && !hasTemplate {
		errors = append(errors, ValidationError{
			Field:   "content/template",
			Message: "either content or template must be provided",
		})
	}

	if hasContent && hasTemplate {
		errors = append(errors, ValidationError{
			Field:   "content/template",
			Message: "content and template cannot be provided simultaneously",
		})
	}

	return errors
}

// validateContentByType validates content based on notification type
func (v *NotificationValidator) validateContentByType(notificationType string, content map[string]interface{}) []ValidationError {
	var errors []ValidationError

	switch notificationType {
	case "email":
		errors = append(errors, v.validateEmailContent(content)...)
	case "slack":
		errors = append(errors, v.validateSlackContent(content)...)
	case "ios_push", "android_push", "in_app":
		errors = append(errors, v.validatePushContent(content)...)
	}

	return errors
}

// validateEmailContent validates email notification content
func (v *NotificationValidator) validateEmailContent(content map[string]interface{}) []ValidationError {
	var errors []ValidationError

	subject, hasSubject := content["subject"].(string)
	if !hasSubject || strings.TrimSpace(subject) == "" {
		errors = append(errors, ValidationError{
			Field:   "content.subject",
			Message: "email subject is required",
		})
	} else if len(subject) > 255 {
		errors = append(errors, ValidationError{
			Field:   "content.subject",
			Message: "email subject cannot exceed 255 characters",
		})
	}

	emailBody, hasEmailBody := content["email_body"].(string)
	if !hasEmailBody || strings.TrimSpace(emailBody) == "" {
		errors = append(errors, ValidationError{
			Field:   "content.email_body",
			Message: "email body is required",
		})
	} else if len(emailBody) > 10000 {
		errors = append(errors, ValidationError{
			Field:   "content.email_body",
			Message: "email body cannot exceed 10000 characters",
		})
	}

	return errors
}

// validateSlackContent validates slack notification content
func (v *NotificationValidator) validateSlackContent(content map[string]interface{}) []ValidationError {
	var errors []ValidationError

	text, hasText := content["text"].(string)
	if !hasText || strings.TrimSpace(text) == "" {
		errors = append(errors, ValidationError{
			Field:   "content.text",
			Message: "slack text is required",
		})
	} else if len(text) > 3000 {
		errors = append(errors, ValidationError{
			Field:   "content.text",
			Message: "slack text cannot exceed 3000 characters",
		})
	}

	return errors
}

// validatePushContent validates push notification content
func (v *NotificationValidator) validatePushContent(content map[string]interface{}) []ValidationError {
	var errors []ValidationError

	title, hasTitle := content["title"].(string)
	if !hasTitle || strings.TrimSpace(title) == "" {
		errors = append(errors, ValidationError{
			Field:   "content.title",
			Message: "push notification title is required",
		})
	} else if len(title) > 255 {
		errors = append(errors, ValidationError{
			Field:   "content.title",
			Message: "push notification title cannot exceed 255 characters",
		})
	}

	body, hasBody := content["body"].(string)
	if !hasBody || strings.TrimSpace(body) == "" {
		errors = append(errors, ValidationError{
			Field:   "content.body",
			Message: "push notification body is required",
		})
	} else if len(body) > 4000 {
		errors = append(errors, ValidationError{
			Field:   "content.body",
			Message: "push notification body cannot exceed 4000 characters",
		})
	}

	return errors
}

// validateTemplate validates template data
func (v *NotificationValidator) validateTemplate(template *models.TemplateData) []ValidationError {
	var errors []ValidationError

	if template.ID == "" {
		errors = append(errors, ValidationError{
			Field:   "template.id",
			Message: "template ID is required",
		})
	} else {
		// Validate UUID format
		uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
		if !uuidRegex.MatchString(strings.ToLower(template.ID)) {
			errors = append(errors, ValidationError{
				Field:   "template.id",
				Message: "template ID must be a valid UUID",
			})
		}
	}

	if template.Data == nil {
		errors = append(errors, ValidationError{
			Field:   "template.data",
			Message: "template data is required",
		})
	} else if len(template.Data) == 0 {
		errors = append(errors, ValidationError{
			Field:   "template.data",
			Message: "template data cannot be empty",
		})
	}

	// Validate version (mandatory field)
	if template.Version <= 0 {
		errors = append(errors, ValidationError{
			Field:   "template.version",
			Message: "template version is required and must be a positive integer",
		})
	}

	return errors
}

// validateFromField validates the from field based on notification type
func (v *NotificationValidator) validateFromField(notificationType string, from *struct {
	Email string `json:"email"`
}) []ValidationError {
	var errors []ValidationError

	if notificationType == "email" {
		if from == nil || strings.TrimSpace(from.Email) == "" {
			errors = append(errors, ValidationError{
				Field:   "from.email",
				Message: "from email is required for email notifications",
			})
		} else {
			// Validate email format
			if _, err := mail.ParseAddress(from.Email); err != nil {
				errors = append(errors, ValidationError{
					Field:   "from.email",
					Message: "invalid email format",
				})
			}

			// Check email length
			if len(from.Email) > 254 {
				errors = append(errors, ValidationError{
					Field:   "from.email",
					Message: "email address cannot exceed 254 characters",
				})
			}
		}
	} else {
		if from != nil {
			errors = append(errors, ValidationError{
				Field:   "from",
				Message: "from field is only allowed for email notifications",
			})
		}
	}

	return errors
}

// validateScheduledAt validates the scheduled_at timestamp
func (v *NotificationValidator) validateScheduledAt(scheduledAt time.Time) []ValidationError {
	var errors []ValidationError

	now := time.Now()

	// Check if scheduled time is in the past
	if scheduledAt.Before(now) {
		errors = append(errors, ValidationError{
			Field:   "scheduled_at",
			Message: "scheduled time cannot be in the past",
		})
	}

	// Check if scheduled time is too far in the future (e.g., 1 year)
	maxScheduledTime := now.AddDate(1, 0, 0)
	if scheduledAt.After(maxScheduledTime) {
		errors = append(errors, ValidationError{
			Field:   "scheduled_at",
			Message: "scheduled time cannot be more than 1 year in the future",
		})
	}

	return errors
}

// ValidateNotificationID validates a notification ID parameter
func (v *NotificationValidator) ValidateNotificationID(notificationID string) ValidationResult {
	var errors []ValidationError

	if notificationID == "" {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "notification ID is required",
		})
		return ValidationResult{IsValid: false, Errors: errors}
	}

	// Validate UUID format
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(strings.ToLower(notificationID)) {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "notification ID must be a valid UUID",
		})
	}

	// Check length constraints
	if len(notificationID) > 36 {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "notification ID cannot exceed 36 characters",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
