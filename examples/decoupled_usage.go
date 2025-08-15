package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager"
	"github.com/gaurav2721/notification-service/notification_manager/templates"
)

// Example demonstrating decoupled template manager usage
func ExampleDecoupledUsage() {
	// Example 1: Using the default template manager
	fmt.Println("=== Example 1: Using Default Template Manager ===")
	useDefaultTemplateManager()

	// Example 2: Using a mock template manager for testing
	fmt.Println("\n=== Example 2: Using Mock Template Manager ===")
	useMockTemplateManager()

	// Example 3: Using a custom template manager implementation
	fmt.Println("\n=== Example 3: Using Custom Template Manager ===")
	useCustomTemplateManager()
}

// useDefaultTemplateManager demonstrates using the default template manager
func useDefaultTemplateManager() {
	// Create notification manager with default template manager
	nm := notification_manager.NewNotificationManagerWithDefaultTemplate(
		nil, // emailService
		nil, // slackService
		nil, // inappService
		nil, // userService
		nil, // scheduler
	)

	// Create a template
	template := &models.Template{
		Name: "Welcome Email",
		Type: models.EmailNotification,
		Content: models.TemplateContent{
			Subject:   "Welcome to our service!",
			EmailBody: "Hello {{name}}, welcome to our platform!",
		},
		RequiredVariables: []string{"name"},
		Description:       "Welcome email template",
	}

	// Create template using the notification manager
	ctx := context.Background()
	response, err := nm.CreateTemplate(ctx, template)
	if err != nil {
		log.Printf("Error creating template: %v", err)
		return
	}

	fmt.Printf("Template created: %+v\n", response)

	// Get predefined templates
	predefinedTemplates := nm.GetPredefinedTemplates()
	fmt.Printf("Found %d predefined templates\n", len(predefinedTemplates))
}

// useMockTemplateManager demonstrates using a mock template manager for testing
func useMockTemplateManager() {
	// Create mock template manager
	mockTemplateManager := templates.NewMockTemplateManager()

	// Create notification manager with mock template manager
	nm := notification_manager.NewNotificationManager(
		nil,                 // emailService
		nil,                 // slackService
		nil,                 // inappService
		nil,                 // userService
		nil,                 // scheduler
		mockTemplateManager, // Using mock instead of real implementation
	)

	// Create a template
	template := &models.Template{
		Name: "Test Template",
		Type: models.EmailNotification,
		Content: models.TemplateContent{
			Subject:   "Test Subject",
			EmailBody: "Hello {{name}}, this is a test!",
		},
		RequiredVariables: []string{"name"},
		Description:       "Test template",
	}

	// Create template
	ctx := context.Background()
	response, err := nm.CreateTemplate(ctx, template)
	if err != nil {
		log.Printf("Error creating template: %v", err)
		return
	}

	fmt.Printf("Mock template created: %+v\n", response)

	// Validate template data
	data := map[string]interface{}{
		"name": "John Doe",
	}

	err = mockTemplateManager.ValidateTemplateData("mock-template-id", data)
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Println("Template data validation successful")
	}
}

// CustomTemplateManager is a custom implementation of the template manager interface
type CustomTemplateManager struct {
	templates map[string]*models.Template
}

// NewCustomTemplateManager creates a new custom template manager
func NewCustomTemplateManager() *CustomTemplateManager {
	return &CustomTemplateManager{
		templates: make(map[string]*models.Template),
	}
}

// Implement the TemplateManager methods
func (c *CustomTemplateManager) CreateTemplate(ctx context.Context, template *models.Template) (*models.TemplateResponse, error) {
	template.ID = "custom-" + template.Name
	c.templates[template.ID] = template

	return &models.TemplateResponse{
		ID:        template.ID,
		Name:      template.Name,
		Type:      string(template.Type),
		Version:   1,
		Status:    "created",
		CreatedAt: template.CreatedAt,
	}, nil
}

func (c *CustomTemplateManager) GetTemplateVersion(ctx context.Context, templateID string, version int) (*models.TemplateVersion, error) {
	template, exists := c.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found")
	}

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

func (c *CustomTemplateManager) GetPredefinedTemplates() []*models.Template {
	return []*models.Template{} // Custom implementation returns empty
}

func (c *CustomTemplateManager) GetTemplateByID(templateID string) (*models.Template, error) {
	template, exists := c.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found")
	}
	return template, nil
}

func (c *CustomTemplateManager) GetTemplateByName(name string) (*models.Template, error) {
	for _, template := range c.templates {
		if template.Name == name {
			return template, nil
		}
	}
	return nil, fmt.Errorf("template not found")
}

func (c *CustomTemplateManager) GetTemplatesByType(templateType models.NotificationType) []*models.Template {
	var filteredTemplates []*models.Template
	for _, template := range c.templates {
		if template.Type == templateType {
			filteredTemplates = append(filteredTemplates, template)
		}
	}
	return filteredTemplates
}

func (c *CustomTemplateManager) ValidateTemplateData(templateID string, data map[string]interface{}) error {
	template, err := c.GetTemplateByID(templateID)
	if err != nil {
		return err
	}

	for _, requiredVar := range template.RequiredVariables {
		if _, exists := data[requiredVar]; !exists {
			return fmt.Errorf("missing required variable: %s", requiredVar)
		}
	}
	return nil
}

func (c *CustomTemplateManager) GetAllTemplates() []*models.Template {
	var allTemplates []*models.Template
	for _, template := range c.templates {
		allTemplates = append(allTemplates, template)
	}
	return allTemplates
}

func (c *CustomTemplateManager) GetTemplateCount() int {
	return len(c.templates)
}

func (c *CustomTemplateManager) GetPredefinedTemplateCount() int {
	return 0 // Custom implementation has no predefined templates
}

// useCustomTemplateManager demonstrates using a custom template manager implementation
func useCustomTemplateManager() {
	// Create custom template manager
	customTemplateManager := NewCustomTemplateManager()

	// Create notification manager with custom template manager
	nm := notification_manager.NewNotificationManager(
		nil,                   // emailService
		nil,                   // slackService
		nil,                   // inappService
		nil,                   // userService
		nil,                   // scheduler
		customTemplateManager, // Using custom implementation
	)

	// Create a template
	template := &models.Template{
		Name: "Custom Template",
		Type: models.SlackNotification,
		Content: models.TemplateContent{
			Text: "Hello {{user}}, you have {{count}} new messages!",
		},
		RequiredVariables: []string{"user", "count"},
		Description:       "Custom Slack notification template",
	}

	// Create template
	ctx := context.Background()
	response, err := nm.CreateTemplate(ctx, template)
	if err != nil {
		log.Printf("Error creating template: %v", err)
		return
	}

	fmt.Printf("Custom template created: %+v\n", response)

	// Get all templates
	allTemplates := customTemplateManager.GetAllTemplates()
	fmt.Printf("Custom template manager has %d templates\n", len(allTemplates))
}
