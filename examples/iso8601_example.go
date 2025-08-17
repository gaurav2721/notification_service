//go:build example
// +build example

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Example demonstrating ISO 8601 timestamp format usage
func ISO8601Example() {
	baseURL := "http://localhost:8080/api/v1"

	// Example 1: Immediate notification (no scheduled_at)
	fmt.Println("=== Example 1: Immediate Notification ===")
	sendImmediateNotification(baseURL)

	// Example 2: Scheduled notification with ISO 8601 timestamp
	fmt.Println("\n=== Example 2: Scheduled Notification with ISO 8601 ===")
	sendScheduledNotification(baseURL)

	// Example 3: Scheduled notification with template
	fmt.Println("\n=== Example 3: Scheduled Notification with Template ===")
	sendScheduledNotificationWithTemplateISO8601(baseURL)
}

func main() {
	ISO8601Example()

	// Also demonstrate different ISO 8601 formats
	demonstrateISO8601Formats()
}

func sendImmediateNotification(baseURL string) {
	notification := map[string]interface{}{
		"type": "email",
		"content": map[string]interface{}{
			"subject":    "Welcome Email",
			"email_body": "Welcome to our platform!",
		},
		"recipients": []string{"user-123", "user-456"},
		// Note: No scheduled_at field means immediate delivery
	}

	sendNotification(baseURL, notification, "Immediate notification")
}

func sendScheduledNotification(baseURL string) {
	// Schedule for 1 hour from now using ISO 8601 format
	scheduledTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	notification := map[string]interface{}{
		"type": "email",
		"content": map[string]interface{}{
			"subject":    "Scheduled Email",
			"email_body": "This email was scheduled for future delivery",
		},
		"recipients":   []string{"user-123", "user-456"},
		"scheduled_at": scheduledTime, // ISO 8601 format: "2024-01-01T12:00:00Z"
	}

	sendNotification(baseURL, notification, "Scheduled notification")
}

func sendScheduledNotificationWithTemplateISO8601(baseURL string) {
	// Schedule for 2 hours from now using ISO 8601 format
	scheduledTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)

	notification := map[string]interface{}{
		"type": "email",
		"template": map[string]interface{}{
			"id": "550e8400-e29b-41d4-a716-446655440000", // Welcome template
			"data": map[string]interface{}{
				"name":            "John Doe",
				"platform":        "Tuskira",
				"username":        "johndoe",
				"email":           "john.doe@example.com",
				"account_type":    "Premium",
				"activation_link": "https://tuskira.com/activate?token=abc123def456",
			},
		},
		"recipients":   []string{"user-789"},
		"scheduled_at": scheduledTime, // ISO 8601 format: "2024-01-01T12:00:00Z"
	}

	sendNotification(baseURL, notification, "Scheduled notification with template")
}

func sendNotification(baseURL string, notification map[string]interface{}, description string) {
	jsonData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling %s: %v", description, err)
		return
	}

	fmt.Printf("Sending %s...\n", description)
	fmt.Printf("Request payload:\n%s\n", string(jsonData))

	resp, err := http.Post(baseURL+"/notifications", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending %s: %v", description, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send %s. Status: %d", description, resp.StatusCode)
		return
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response for %s: %v", description, err)
		return
	}

	fmt.Printf("✅ %s sent successfully:\n", description)
	fmt.Printf("  ID: %s\n", response["id"])
	fmt.Printf("  Status: %s\n", response["status"])
	fmt.Printf("  Message: %s\n", response["message"])
	if response["scheduled_at"] != nil {
		fmt.Printf("  Scheduled At: %s\n", response["scheduled_at"])
	}
	fmt.Println()
}

// Helper function to demonstrate different ISO 8601 formats
func demonstrateISO8601Formats() {
	fmt.Println("=== ISO 8601 Timestamp Formats ===")

	now := time.Now()

	// RFC3339 format (most common)
	fmt.Printf("RFC3339: %s\n", now.Format(time.RFC3339))

	// RFC3339Nano format (with nanoseconds)
	fmt.Printf("RFC3339Nano: %s\n", now.Format(time.RFC3339Nano))

	// Custom format with timezone offset
	fmt.Printf("Custom with offset: %s\n", now.Format("2006-01-02T15:04:05-07:00"))

	// UTC format
	fmt.Printf("UTC: %s\n", now.UTC().Format(time.RFC3339))

	fmt.Println()
}
