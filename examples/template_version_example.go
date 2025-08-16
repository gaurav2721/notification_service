package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// TemplateVersionExample demonstrates the new version field functionality
func TemplateVersionExample() {
	baseURL := "http://localhost:8080/api/v1"

	fmt.Println("=== Template Version Example ===")
	fmt.Println()

	// Example 1: Send notification with specific version
	fmt.Println("1. Sending email notification with specific template version...")
	sendNotificationWithVersion(baseURL, "550e8400-e29b-41d4-a716-446655440000", 1)

	fmt.Println()

	// Example 2: Send notification with different version
	fmt.Println("2. Sending email notification with different template version...")
	sendNotificationWithVersion(baseURL, "550e8400-e29b-41d4-a716-446655440000", 2)

	fmt.Println()

	// Example 3: Demonstrate validation error when version is missing
	fmt.Println("3. Demonstrating validation error when version is missing...")
	sendNotificationWithMissingVersion(baseURL, "550e8400-e29b-41d4-a716-446655440000")

	fmt.Println()

	// Example 4: Schedule notification with specific version
	fmt.Println("4. Scheduling email notification with specific template version...")
	scheduleNotificationWithVersion(baseURL, "550e8400-e29b-41d4-a716-446655440000", 1)

	fmt.Println()

	// Example 5: Send Slack notification with version
	fmt.Println("5. Sending Slack notification with specific template version...")
	sendSlackNotificationWithVersion(baseURL, "550e8400-e29b-41d4-a716-446655440003", 1)

	fmt.Println()

	// Example 6: Send In-App notification with version
	fmt.Println("6. Sending In-App notification with specific template version...")
	sendInAppNotificationWithVersion(baseURL, "550e8400-e29b-41d4-a716-446655440005", 1)
}

func sendNotificationWithVersion(baseURL, templateID string, version int) {
	url := fmt.Sprintf("%s/notifications", baseURL)

	payload := map[string]interface{}{
		"type": "email",
		"template": map[string]interface{}{
			"id":      templateID,
			"version": version,
			"data": map[string]interface{}{
				"name":            "John Doe",
				"platform":        "Tuskira",
				"username":        "johndoe",
				"email":           "john.doe@example.com",
				"account_type":    "Premium",
				"activation_link": "https://tuskira.com/activate?token=abc123def456",
			},
		},
		"recipients": []string{"user-001"},
		"from": map[string]interface{}{
			"email": "noreply@company.com",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending notification with version: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Successfully sent notification with template version %d\n", version)
	} else {
		fmt.Printf("❌ Failed to send notification with version. Status: %d\n", resp.StatusCode)
	}
}

func scheduleNotificationWithVersion(baseURL, templateID string, version int) {
	url := fmt.Sprintf("%s/notifications", baseURL)

	// Schedule for 1 hour from now
	scheduledTime := time.Now().Add(time.Hour)

	payload := map[string]interface{}{
		"type": "email",
		"template": map[string]interface{}{
			"id":      templateID,
			"version": version,
			"data": map[string]interface{}{
				"name":            "Bob Wilson",
				"platform":        "Tuskira",
				"username":        "bobwilson",
				"email":           "bob.wilson@example.com",
				"account_type":    "Basic",
				"activation_link": "https://tuskira.com/activate?token=ghi789jkl012",
			},
		},
		"recipients":   []string{"user-003"},
		"scheduled_at": scheduledTime.Format(time.RFC3339),
		"from": map[string]interface{}{
			"email": "noreply@company.com",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error scheduling notification with version: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Successfully scheduled notification with template version %d\n", version)
	} else {
		fmt.Printf("❌ Failed to schedule notification with version. Status: %d\n", resp.StatusCode)
	}
}

func sendSlackNotificationWithVersion(baseURL, templateID string, version int) {
	url := fmt.Sprintf("%s/notifications", baseURL)

	payload := map[string]interface{}{
		"type": "slack",
		"template": map[string]interface{}{
			"id":      templateID,
			"version": version,
			"data": map[string]interface{}{
				"alert_type":        "System Health",
				"system_name":       "Notification Service",
				"severity":          "Info",
				"environment":       "Development",
				"message":           "Template versioning feature is working correctly",
				"timestamp":         time.Now().Format(time.RFC3339),
				"action_required":   "No action required",
				"affected_services": "Template processing",
				"dashboard_link":    "https://dashboard.example.com/health",
			},
		},
		"recipients": []string{"user-001", "user-002"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending Slack notification with version: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Successfully sent Slack notification with template version %d\n", version)
	} else {
		fmt.Printf("❌ Failed to send Slack notification with version. Status: %d\n", resp.StatusCode)
	}
}

func sendInAppNotificationWithVersion(baseURL, templateID string, version int) {
	url := fmt.Sprintf("%s/notifications", baseURL)

	payload := map[string]interface{}{
		"type": "in_app",
		"template": map[string]interface{}{
			"id":      templateID,
			"version": version,
			"data": map[string]interface{}{
				"order_id":       "ORD-2024-VERSION-TEST",
				"status":         "Processing",
				"item_count":     1,
				"total_amount":   "199.99",
				"status_message": "Your order is being processed with template versioning",
				"action_button":  "View Order",
			},
		},
		"recipients": []string{"user-001"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending In-App notification with version: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Successfully sent In-App notification with template version %d\n", version)
	} else {
		fmt.Printf("❌ Failed to send In-App notification with version. Status: %d\n", resp.StatusCode)
	}
}

func sendNotificationWithMissingVersion(baseURL, templateID string) {
	url := fmt.Sprintf("%s/notifications", baseURL)

	// This payload is missing the version field, which should cause a validation error
	payload := map[string]interface{}{
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
		"recipients": []string{"user-002"},
		"from": map[string]interface{}{
			"email": "noreply@company.com",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending notification with missing version: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("✅ Correctly received validation error for missing version field")
	} else {
		fmt.Printf("❌ Expected validation error but got status: %d\n", resp.StatusCode)
	}
}

// To run this example, call TemplateVersionExample() from your main function
// or create a separate main function in a different package
