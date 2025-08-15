package templates

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/google/uuid"
)

// TemplateManager handles all template-related operations
type TemplateManager struct {
	templates     map[string]*models.Template
	templateMutex sync.RWMutex
	initialized   bool
}

// NewTemplateManager creates a new template manager instance
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates:     make(map[string]*models.Template),
		templateMutex: sync.RWMutex{},
		initialized:   false,
	}

	// Load predefined templates on startup
	tm.loadPredefinedTemplates()

	return tm
}

// loadPredefinedTemplates loads all predefined templates into the manager
func (tm *TemplateManager) loadPredefinedTemplates() {
	tm.templateMutex.Lock()
	defer tm.templateMutex.Unlock()

	predefinedTemplates := models.PredefinedTemplates()

	for _, template := range predefinedTemplates {
		tm.templates[template.ID] = template
		log.Printf("Loaded predefined template: %s (ID: %s)", template.Name, template.ID)
	}

	tm.initialized = true
	log.Printf("Loaded %d predefined templates", len(predefinedTemplates))
}

// CreateTemplate creates a new notification template
func (tm *TemplateManager) CreateTemplate(ctx context.Context, template *models.Template) (*models.TemplateResponse, error) {
	tm.templateMutex.Lock()
	defer tm.templateMutex.Unlock()

	// Validate template content
	if err := template.Content.ValidateTemplateContent(template.Type); err != nil {
		return nil, err
	}

	// Generate ID if not provided
	if template.ID == "" {
		template.ID = uuid.New().String()
	}

	// Set version and status
	template.Version = 1
	template.Status = "created"
	template.CreatedAt = time.Now()

	// Store template
	tm.templates[template.ID] = template

	// Return response
	return &models.TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      string(template.Type),
		Version:   template.Version,
		Status:    template.Status,
		CreatedAt: template.CreatedAt,
	}, nil
}

// GetTemplateVersion retrieves a specific version of a notification template
func (tm *TemplateManager) GetTemplateVersion(ctx context.Context, templateID string, version int) (*models.TemplateVersion, error) {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	template, exists := tm.templates[templateID]
	if !exists {
		return nil, models.ErrTemplateNotFound
	}

	// For now, we only support version 1
	// In a real implementation, you would store multiple versions
	if version != 1 {
		return nil, models.ErrTemplateNotFound
	}

	// Return template version
	return &models.TemplateVersion{
		ID:                template.ID,
		Name:              template.Name,
		Type:              template.Type,
		Version:           template.Version,
		Content:           template.Content,
		RequiredVariables: template.RequiredVariables,
		Description:       template.Description,
		Status:            template.Status,
		CreatedAt:         template.CreatedAt,
	}, nil
}

// GetPredefinedTemplates returns all predefined templates
func (tm *TemplateManager) GetPredefinedTemplates() []*models.Template {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	var predefinedTemplates []*models.Template
	for _, template := range tm.templates {
		// Check if it's a predefined template by looking at the fixed IDs
		if isPredefinedTemplateID(template.ID) {
			predefinedTemplates = append(predefinedTemplates, template)
		}
	}

	return predefinedTemplates
}

// GetTemplateByID returns a template by ID
func (tm *TemplateManager) GetTemplateByID(templateID string) (*models.Template, error) {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	template, exists := tm.templates[templateID]
	if !exists {
		return nil, models.ErrTemplateNotFound
	}

	return template, nil
}

// GetTemplateByName returns a template by name
func (tm *TemplateManager) GetTemplateByName(name string) (*models.Template, error) {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	for _, template := range tm.templates {
		if template.Name == name {
			return template, nil
		}
	}

	return nil, models.ErrTemplateNotFound
}

// GetTemplatesByType returns all templates of a specific type
func (tm *TemplateManager) GetTemplatesByType(templateType models.NotificationType) []*models.Template {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	var filteredTemplates []*models.Template
	for _, template := range tm.templates {
		if template.Type == templateType {
			filteredTemplates = append(filteredTemplates, template)
		}
	}

	return filteredTemplates
}

// ValidateTemplateData validates that the provided data contains all required variables for a template
func (tm *TemplateManager) ValidateTemplateData(templateID string, data map[string]interface{}) error {
	template, err := tm.GetTemplateByID(templateID)
	if err != nil {
		return err
	}

	return template.ValidateRequiredVariables(data)
}

// GetAllTemplates returns all templates (both predefined and custom)
func (tm *TemplateManager) GetAllTemplates() []*models.Template {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	var allTemplates []*models.Template
	for _, template := range tm.templates {
		allTemplates = append(allTemplates, template)
	}

	return allTemplates
}

// GetTemplateCount returns the total number of templates
func (tm *TemplateManager) GetTemplateCount() int {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	return len(tm.templates)
}

// GetPredefinedTemplateCount returns the number of predefined templates
func (tm *TemplateManager) GetPredefinedTemplateCount() int {
	tm.templateMutex.RLock()
	defer tm.templateMutex.RUnlock()

	count := 0
	for _, template := range tm.templates {
		if isPredefinedTemplateID(template.ID) {
			count++
		}
	}

	return count
}

// isPredefinedTemplateID checks if a template ID is one of the predefined ones
func isPredefinedTemplateID(templateID string) bool {
	predefinedIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000", // Welcome Email
		"550e8400-e29b-41d4-a716-446655440001", // Password Reset
		"550e8400-e29b-41d4-a716-446655440002", // Order Confirmation
		"550e8400-e29b-41d4-a716-446655440003", // System Alert
		"550e8400-e29b-41d4-a716-446655440004", // Deployment Notification
		"550e8400-e29b-41d4-a716-446655440005", // Order Status Update
		"550e8400-e29b-41d4-a716-446655440006", // Payment Reminder
	}

	for _, id := range predefinedIDs {
		if templateID == id {
			return true
		}
	}
	return false
}
