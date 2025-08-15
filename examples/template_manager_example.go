package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager/templates"
)

// Example demonstrating the new TemplateManager

func templateManagerExample() {
	fmt.Println("=== Template Manager Example ===")

	// Create a new template manager
	templateManager := templates.NewTemplateManager()

	// Example 1: Get all predefined templates
	fmt.Println("\n1. Getting all predefined templates:")
	predefinedTemplates := templateManager.GetPredefinedTemplates()
	for _, template := range predefinedTemplates {
		fmt.Printf("  - %s (ID: %s, Type: %s)\n", template.Name, template.ID, template.Type)
	}

	// Example 2: Get templates by type
	fmt.Println("\n2. Getting email templates:")
	emailTemplates := templateManager.GetTemplatesByType(models.EmailNotification)
	for _, template := range emailTemplates {
		fmt.Printf("  - %s: %s\n", template.Name, template.Description)
	}

	// Example 3: Get a specific template by ID
	fmt.Println("\n3. Getting Welcome Email Template by ID:")
	welcomeTemplate, err := templateManager.GetTemplateByID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		log.Printf("Error getting template: %v", err)
	} else {
		fmt.Printf("  Template: %s\n", welcomeTemplate.Name)
		fmt.Printf("  Required Variables: %v\n", welcomeTemplate.RequiredVariables)
	}

	// Example 4: Get a template by name
	fmt.Println("\n4. Getting System Alert Template by name:")
	systemAlertTemplate, err := templateManager.GetTemplateByName("System Alert Template")
	if err != nil {
		log.Printf("Error getting template: %v", err)
	} else {
		fmt.Printf("  Template: %s (ID: %s)\n", systemAlertTemplate.Name, systemAlertTemplate.ID)
		fmt.Printf("  Type: %s\n", systemAlertTemplate.Type)
	}

	// Example 5: Validate template data
	fmt.Println("\n5. Validating template data:")
	testData := map[string]interface{}{
		"name":            "John Doe",
		"platform":        "Tuskira",
		"username":        "johndoe",
		"email":           "john.doe@example.com",
		"account_type":    "Premium",
		"activation_link": "https://tuskira.com/activate?token=abc123",
	}

	err = templateManager.ValidateTemplateData("550e8400-e29b-41d4-a716-446655440000", testData)
	if err != nil {
		fmt.Printf("  Validation failed: %v\n", err)
	} else {
		fmt.Printf("  Validation passed for Welcome Email Template\n")
	}

	// Example 6: Test validation with missing data
	fmt.Println("\n6. Testing validation with missing data:")
	incompleteData := map[string]interface{}{
		"name":     "John Doe",
		"platform": "Tuskira",
		// Missing required variables
	}

	err = templateManager.ValidateTemplateData("550e8400-e29b-41d4-a716-446655440000", incompleteData)
	if err != nil {
		fmt.Printf("  Validation correctly failed: %v\n", err)
	} else {
		fmt.Printf("  Validation unexpectedly passed\n")
	}

	// Example 7: Create a custom template
	fmt.Println("\n7. Creating a custom template:")
	customTemplate := &models.Template{
		Name: "Custom Welcome Template",
		Type: models.EmailNotification,
		Content: models.TemplateContent{
			Subject:   "Welcome to {{company}}, {{name}}!",
			EmailBody: "Hi {{name}},\n\nWelcome to {{company}}! We're excited to have you join our team.\n\nBest regards,\nThe {{company}} Team",
		},
		RequiredVariables: []string{"name", "company"},
		Description:       "Custom welcome email template",
	}

	ctx := context.Background()
	response, err := templateManager.CreateTemplate(ctx, customTemplate)
	if err != nil {
		log.Printf("Error creating template: %v", err)
	} else {
		fmt.Printf("  Created custom template: %s (ID: %s)\n", response.Name, response.ID)
	}

	// Example 8: Get template statistics
	fmt.Println("\n8. Template statistics:")
	totalTemplates := templateManager.GetTemplateCount()
	predefinedCount := templateManager.GetPredefinedTemplateCount()
	customCount := totalTemplates - predefinedCount

	fmt.Printf("  Total templates: %d\n", totalTemplates)
	fmt.Printf("  Predefined templates: %d\n", predefinedCount)
	fmt.Printf("  Custom templates: %d\n", customCount)

	// Example 9: Get all templates
	fmt.Println("\n9. Getting all templates:")
	allTemplates := templateManager.GetAllTemplates()
	for _, template := range allTemplates {
		templateType := "Custom"
		if isPredefinedTemplateID(template.ID) {
			templateType = "Predefined"
		}
		fmt.Printf("  - %s (%s): %s\n", template.Name, templateType, template.Description)
	}
}

// Helper function to check if a template ID is predefined
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

// Example of using predefined templates in notifications
func predefinedTemplatesUsageExample() {
	fmt.Println("\n=== Predefined Templates Usage Example ===")

	templateManager := templates.NewTemplateManager()

	// Example: Using Welcome Email Template
	fmt.Println("\n1. Using Welcome Email Template:")
	welcomeData := map[string]interface{}{
		"name":            "Jane Smith",
		"platform":        "Tuskira",
		"username":        "janesmith",
		"email":           "jane.smith@example.com",
		"account_type":    "Standard",
		"activation_link": "https://tuskira.com/activate?token=def456ghi789",
	}

	err := templateManager.ValidateTemplateData("550e8400-e29b-41d4-a716-446655440000", welcomeData)
	if err != nil {
		fmt.Printf("  Validation failed: %v\n", err)
	} else {
		fmt.Printf("  Welcome email template data validated successfully\n")
		fmt.Printf("  Ready to send welcome email to: %s\n", welcomeData["name"])
	}

	// Example: Using System Alert Template
	fmt.Println("\n2. Using System Alert Template:")
	alertData := map[string]interface{}{
		"alert_type":        "Database Connection",
		"system_name":       "User Service",
		"severity":          "Critical",
		"environment":       "Production",
		"message":           "Database connection timeout after 30 seconds",
		"timestamp":         "2024-01-01T10:00:00Z",
		"action_required":   "Check database connectivity and restart service if needed",
		"affected_services": "User authentication, profile management",
		"dashboard_link":    "https://grafana.company.com/d/user-service",
	}

	err = templateManager.ValidateTemplateData("550e8400-e29b-41d4-a716-446655440003", alertData)
	if err != nil {
		fmt.Printf("  Validation failed: %v\n", err)
	} else {
		fmt.Printf("  System alert template data validated successfully\n")
		fmt.Printf("  Ready to send alert for: %s\n", alertData["system_name"])
	}

	// Example: Using Order Status Update Template
	fmt.Println("\n3. Using Order Status Update Template:")
	orderData := map[string]interface{}{
		"order_id":       "ORD-2024-001",
		"status":         "Shipped",
		"item_count":     "3",
		"total_amount":   "299.99",
		"status_message": "Your order has been shipped and is on its way!",
		"action_button":  "Track Order",
	}

	err = templateManager.ValidateTemplateData("550e8400-e29b-41d4-a716-446655440005", orderData)
	if err != nil {
		fmt.Printf("  Validation failed: %v\n", err)
	} else {
		fmt.Printf("  Order status template data validated successfully\n")
		fmt.Printf("  Ready to send order update for: %s\n", orderData["order_id"])
	}
}
