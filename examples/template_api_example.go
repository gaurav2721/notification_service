package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Example usage of the Template APIs

func templateAPIExample() {
	baseURL := "http://localhost:8080/api/v1"

	// Example 1: Create a Welcome Email Template
	fmt.Println("=== Creating Welcome Email Template ===")
	welcomeTemplate := map[string]interface{}{
		"name": "Welcome Email Template",
		"type": "email",
		"content": map[string]interface{}{
			"subject":    "Welcome to {{platform}}, {{name}}!",
			"email_body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team",
		},
		"required_variables": []string{"name", "platform", "username", "email", "account_type", "activation_link"},
		"description":        "Welcome email template for new user onboarding",
	}

	templateID := createTemplate(baseURL, welcomeTemplate)
	if templateID == "" {
		log.Fatal("Failed to create template")
	}

	// Example 2: Get Template Version
	fmt.Println("\n=== Getting Template Version ===")
	getTemplateVersion(baseURL, templateID, 1)

	// Example 3: Send Notification Using Template
	fmt.Println("\n=== Sending Notification Using Template ===")
	sendNotificationWithTemplate(baseURL, templateID)

	// Example 4: Create Slack Alert Template
	fmt.Println("\n=== Creating Slack Alert Template ===")
	slackTemplate := map[string]interface{}{
		"name": "System Alert Template",
		"type": "slack",
		"content": map[string]interface{}{
			"text": "🚨 *{{alert_type}} Alert*\n\n*System:* {{system_name}}\n*Severity:* {{severity}}\n*Environment:* {{environment}}\n*Message:* {{message}}\n*Timestamp:* {{timestamp}}\n*Action Required:* {{action_required}}\n\n*Affected Services:* {{affected_services}}\n*Dashboard:* {{dashboard_link}}\n\nPlease take immediate action if this is a critical alert.",
		},
		"required_variables": []string{"alert_type", "system_name", "severity", "environment", "message", "timestamp", "action_required", "affected_services", "dashboard_link"},
		"description":        "Slack alert template for system monitoring",
	}

	slackTemplateID := createTemplate(baseURL, slackTemplate)
	if slackTemplateID != "" {
		fmt.Printf("Slack template created with ID: %s\n", slackTemplateID)
	}

	// Example 5: Create In-App Notification Template
	fmt.Println("\n=== Creating In-App Notification Template ===")
	inAppTemplate := map[string]interface{}{
		"name": "Order Status Update Template",
		"type": "in_app",
		"content": map[string]interface{}{
			"title": "Order #{{order_id}} - {{status}}",
			"body":  "Your order has been {{status}}.\n\n*Order Details:*\n- Items: {{item_count}} items\n- Total: ${{total_amount}}\n- Status: {{status}}\n\n{{status_message}}\n\n{{action_button}}",
		},
		"required_variables": []string{"order_id", "status", "item_count", "total_amount", "status_message", "action_button"},
		"description":        "In-app notification template for order status updates",
	}

	inAppTemplateID := createTemplate(baseURL, inAppTemplate)
	if inAppTemplateID != "" {
		fmt.Printf("In-app template created with ID: %s\n", inAppTemplateID)
	}
}

func createTemplate(baseURL string, template map[string]interface{}) string {
	jsonData, err := json.Marshal(template)
	if err != nil {
		log.Printf("Error marshaling template: %v", err)
		return ""
	}

	resp, err := http.Post(baseURL+"/templates", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating template: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to create template. Status: %d", resp.StatusCode)
		return ""
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return ""
	}

	templateID, ok := response["id"].(string)
	if !ok {
		log.Printf("Invalid response format")
		return ""
	}

	fmt.Printf("Template created successfully with ID: %s\n", templateID)
	fmt.Printf("Template Name: %s\n", response["name"])
	fmt.Printf("Template Type: %s\n", response["type"])
	fmt.Printf("Template Version: %v\n", response["version"])
	fmt.Printf("Template Status: %s\n", response["status"])

	return templateID
}

func getTemplateVersion(baseURL, templateID string, version int) {
	url := fmt.Sprintf("%s/templates/%s/versions/%d", baseURL, templateID, version)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting template version: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get template version. Status: %d", resp.StatusCode)
		return
	}

	var template map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&template); err != nil {
		log.Printf("Error decoding template: %v", err)
		return
	}

	fmt.Printf("Template Version Retrieved:\n")
	fmt.Printf("  ID: %s\n", template["id"])
	fmt.Printf("  Name: %s\n", template["name"])
	fmt.Printf("  Type: %s\n", template["type"])
	fmt.Printf("  Version: %v\n", template["version"])
	fmt.Printf("  Status: %s\n", template["status"])
	fmt.Printf("  Created At: %s\n", template["created_at"])
}

func sendNotificationWithTemplate(baseURL, templateID string) {
	notification := map[string]interface{}{
		"type": "email",
		"template": map[string]interface{}{
			"id": templateID,
			"data": map[string]interface{}{
				"name":            "John Doe",
				"platform":        "Tuskira",
				"username":        "johndoe",
				"email":           "john.doe@example.com",
				"account_type":    "Premium",
				"activation_link": "https://tuskira.com/activate?token=abc123def456",
			},
		},
		"recipients": []string{"user-123", "user-456"},
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	resp, err := http.Post(baseURL+"/notifications", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send notification. Status: %d", resp.StatusCode)
		return
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return
	}

	fmt.Printf("Notification sent successfully:\n")
	fmt.Printf("  ID: %s\n", response["id"])
	fmt.Printf("  Status: %s\n", response["status"])
	fmt.Printf("  Message: %s\n", response["message"])
}

// Example of scheduled notification with template
func sendScheduledNotificationWithTemplate(baseURL, templateID string) {
	scheduledTime := time.Now().Add(1 * time.Hour)

	notification := map[string]interface{}{
		"type": "email",
		"template": map[string]interface{}{
			"id": templateID,
			"data": map[string]interface{}{
				"name":            "Jane Smith",
				"platform":        "Tuskira",
				"username":        "janesmith",
				"email":           "jane.smith@example.com",
				"account_type":    "Standard",
				"activation_link": "https://tuskira.com/activate?token=def456ghi789",
			},
		},
		"recipients":   []string{"user-789"},
		"scheduled_at": scheduledTime,
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling scheduled notification: %v", err)
		return
	}

	resp, err := http.Post(baseURL+"/notifications", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error scheduling notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to schedule notification. Status: %d", resp.StatusCode)
		return
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return
	}

	fmt.Printf("Notification scheduled successfully:\n")
	fmt.Printf("  ID: %s\n", response["id"])
	fmt.Printf("  Status: %s\n", response["status"])
	fmt.Printf("  Message: %s\n", response["message"])
	fmt.Printf("  Scheduled At: %v\n", response["scheduled_at"])
}
