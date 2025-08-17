package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gaurav2721/notification-service/models"
)

// TemplateValidator provides validation methods for template requests
type TemplateValidator struct{}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{}
}

// ValidateTemplateRequest validates a complete template creation request
func (v *TemplateValidator) ValidateTemplateRequest(request *models.TemplateRequest) ValidationResult {
	var errors []ValidationError

	// Validate basic required fields
	if nameErrors := v.validateTemplateName(request.Name); len(nameErrors) > 0 {
		errors = append(errors, nameErrors...)
	}

	if typeErrors := v.validateTemplateType(request.Type); len(typeErrors) > 0 {
		errors = append(errors, typeErrors...)
	}

	if contentErrors := v.validateTemplateContent(request.Content, request.Type); len(contentErrors) > 0 {
		errors = append(errors, contentErrors...)
	}

	if variableErrors := v.validateRequiredVariables(request.RequiredVariables); len(variableErrors) > 0 {
		errors = append(errors, variableErrors...)
	}

	if descriptionErrors := v.validateTemplateDescription(request.Description); len(descriptionErrors) > 0 {
		errors = append(errors, descriptionErrors...)
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateTemplateData validates template data for usage
func (v *TemplateValidator) ValidateTemplateData(templateData *models.TemplateData) ValidationResult {
	var errors []ValidationError

	if templateData == nil {
		errors = append(errors, ValidationError{
			Field:   "template",
			Message: "template data cannot be nil",
		})
		return ValidationResult{IsValid: false, Errors: errors}
	}

	if templateData.ID == "" {
		errors = append(errors, ValidationError{
			Field:   "template.id",
			Message: "template ID is required",
		})
	}

	if templateData.Version <= 0 {
		errors = append(errors, ValidationError{
			Field:   "template.version",
			Message: "template version must be greater than 0",
		})
	}

	if templateData.Data == nil {
		errors = append(errors, ValidationError{
			Field:   "template.data",
			Message: "template data cannot be nil",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateTemplateID validates a template ID parameter
func (v *TemplateValidator) ValidateTemplateID(templateID string) ValidationResult {
	var errors []ValidationError

	if templateID == "" {
		errors = append(errors, ValidationError{
			Field:   "templateId",
			Message: "template ID is required",
		})
		return ValidationResult{IsValid: false, Errors: errors}
	}

	// Validate UUID format
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(templateID) {
		errors = append(errors, ValidationError{
			Field:   "templateId",
			Message: "template ID must be a valid UUID",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateTemplateVersion validates a template version parameter
func (v *TemplateValidator) ValidateTemplateVersion(version int) ValidationResult {
	var errors []ValidationError

	if version <= 0 {
		errors = append(errors, ValidationError{
			Field:   "version",
			Message: "template version must be greater than 0",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// validateTemplateName validates the template name
func (v *TemplateValidator) validateTemplateName(name string) []ValidationError {
	var errors []ValidationError

	if name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "template name is required",
		})
		return errors
	}

	if len(strings.TrimSpace(name)) == 0 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "template name cannot be empty or whitespace only",
		})
	}

	if len(name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "template name cannot exceed 100 characters",
		})
	}

	// Validate name format (alphanumeric, spaces, hyphens, underscores)
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !nameRegex.MatchString(name) {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "template name can only contain alphanumeric characters, spaces, hyphens, and underscores",
		})
	}

	return errors
}

// validateTemplateType validates the template type
func (v *TemplateValidator) validateTemplateType(templateType models.NotificationType) []ValidationError {
	var errors []ValidationError

	if templateType == "" {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: "template type is required",
		})
		return errors
	}

	validTypes := map[models.NotificationType]bool{
		models.EmailNotification: true,
		models.SlackNotification: true,
		models.InAppNotification: true,
	}

	if !validTypes[templateType] {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("invalid template type. Must be one of: %s", getValidTemplateTypes()),
		})
	}

	return errors
}

// validateTemplateContent validates the template content based on type
func (v *TemplateValidator) validateTemplateContent(content models.TemplateContent, templateType models.NotificationType) []ValidationError {
	var errors []ValidationError

	switch templateType {
	case models.EmailNotification:
		if content.Subject == "" {
			errors = append(errors, ValidationError{
				Field:   "content.subject",
				Message: "email template subject is required",
			})
		}
		if content.EmailBody == "" {
			errors = append(errors, ValidationError{
				Field:   "content.email_body",
				Message: "email template body is required",
			})
		}
		if len(content.Subject) > 200 {
			errors = append(errors, ValidationError{
				Field:   "content.subject",
				Message: "email subject cannot exceed 200 characters",
			})
		}
		if len(content.EmailBody) > 10000 {
			errors = append(errors, ValidationError{
				Field:   "content.email_body",
				Message: "email body cannot exceed 10000 characters",
			})
		}

	case models.SlackNotification:
		if content.Text == "" {
			errors = append(errors, ValidationError{
				Field:   "content.text",
				Message: "slack template text is required",
			})
		}
		if len(content.Text) > 3000 {
			errors = append(errors, ValidationError{
				Field:   "content.text",
				Message: "slack text cannot exceed 3000 characters",
			})
		}

	case models.InAppNotification:
		if content.Title == "" {
			errors = append(errors, ValidationError{
				Field:   "content.title",
				Message: "in-app template title is required",
			})
		}
		if content.Body == "" {
			errors = append(errors, ValidationError{
				Field:   "content.body",
				Message: "in-app template body is required",
			})
		}
		if len(content.Title) > 100 {
			errors = append(errors, ValidationError{
				Field:   "content.title",
				Message: "in-app title cannot exceed 100 characters",
			})
		}
		if len(content.Body) > 500 {
			errors = append(errors, ValidationError{
				Field:   "content.body",
				Message: "in-app body cannot exceed 500 characters",
			})
		}
	}

	return errors
}

// validateRequiredVariables validates the required variables list
func (v *TemplateValidator) validateRequiredVariables(variables []string) []ValidationError {
	var errors []ValidationError

	if variables == nil {
		errors = append(errors, ValidationError{
			Field:   "required_variables",
			Message: "required variables list cannot be nil",
		})
		return errors
	}

	if len(variables) == 0 {
		errors = append(errors, ValidationError{
			Field:   "required_variables",
			Message: "at least one required variable must be specified",
		})
		return errors
	}

	// Check for duplicate variables
	seen := make(map[string]bool)
	for i, variable := range variables {
		if variable == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("required_variables[%d]", i),
				Message: "variable name cannot be empty",
			})
			continue
		}

		if seen[variable] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("required_variables[%d]", i),
				Message: fmt.Sprintf("duplicate variable name: %s", variable),
			})
			continue
		}

		// Validate variable name format (alphanumeric and underscores only)
		varRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
		if !varRegex.MatchString(variable) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("required_variables[%d]", i),
				Message: fmt.Sprintf("invalid variable name format: %s. Must start with letter or underscore and contain only alphanumeric characters and underscores", variable),
			})
		}

		seen[variable] = true
	}

	return errors
}

// validateTemplateDescription validates the template description
func (v *TemplateValidator) validateTemplateDescription(description string) []ValidationError {
	var errors []ValidationError

	if description != "" && len(description) > 500 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "template description cannot exceed 500 characters",
		})
	}

	return errors
}

// getValidTemplateTypes returns a comma-separated list of valid template types
func getValidTemplateTypes() string {
	return "email, slack, in_app"
}
