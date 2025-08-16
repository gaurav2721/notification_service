package templates

import (
	"context"

	"github.com/gaurav2721/notification-service/models"
)

// TemplateManager defines the contract for template management operations
type TemplateManager interface {
	// CreateTemplate creates a new notification template
	CreateTemplate(ctx context.Context, template *models.Template) (*models.TemplateResponse, error)

	// GetTemplateVersion retrieves a specific version of a notification template
	GetTemplateVersion(ctx context.Context, templateID string, version int) (*models.TemplateVersion, error)

	// GetPredefinedTemplates returns all predefined templates
	GetPredefinedTemplates() []*models.Template

	// GetTemplateByID returns a template by ID (latest version)
	GetTemplateByID(templateID string) (*models.Template, error)

	// GetTemplateByIDAndVersion returns a specific version of a template
	GetTemplateByIDAndVersion(templateID string, version int) (*models.Template, error)

	// GetTemplateByName returns a template by name
	GetTemplateByName(name string) (*models.Template, error)

	// GetTemplatesByType returns all templates of a specific type
	GetTemplatesByType(templateType models.NotificationType) []*models.Template

	// ValidateTemplateData validates that the provided data contains all required variables for a template
	ValidateTemplateData(templateID string, data map[string]interface{}) error

	// GetAllTemplates returns all templates (both predefined and custom)
	GetAllTemplates() []*models.Template

	// GetTemplateCount returns the total number of templates
	GetTemplateCount() int

	// GetPredefinedTemplateCount returns the number of predefined templates
	GetPredefinedTemplateCount() int
}
