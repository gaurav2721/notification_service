package templates

import (
	"github.com/gaurav2721/notification-service/models"
)

// TemplateManager defines the contract for template management operations
type TemplateManager interface {
	// CreateTemplate creates a new notification template
	CreateTemplate(template *models.Template) (*models.TemplateResponse, error)

	// GetTemplateVersion retrieves a specific version of a notification template
	GetTemplateVersion(templateID string, version int) (*models.TemplateVersion, error)

	// GetPredefinedTemplates returns all predefined templates
	GetPredefinedTemplates() []*models.Template

	// GetTemplateByIDAndVersion returns a specific version of a template
	GetTemplateByIDAndVersion(templateID string, version int) (*models.Template, error)
}
