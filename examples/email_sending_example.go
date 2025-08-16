package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gaurav2721/notification-service/external_services/consumers"
	"github.com/gaurav2721/notification-service/models"
)

// Example demonstrating the complete email sending flow
func main() {
	fmt.Println("=== Email Notification Service Example ===")
	fmt.Println("This example demonstrates how EmailNotificationRequest payloads are processed")
	fmt.Println("and emails are sent through the notification service.\n")

	// Step 1: Create an EmailNotificationRequest
	fmt.Println("Step 1: Creating EmailNotificationRequest...")
	emailNotification := createEmailNotificationRequest()

	// Step 2: Convert to JSON payload (simulating what would be sent to Kafka)
	fmt.Println("Step 2: Converting to JSON payload...")
	payload := convertToJSONPayload(emailNotification)

	// Step 3: Create notification message (simulating Kafka message)
	fmt.Println("Step 3: Creating notification message...")
	message := createNotificationMessage(payload)

	// Step 4: Process the notification using email processor
	fmt.Println("Step 4: Processing notification with email processor...")
	processEmailNotification(message)

	fmt.Println("\n=== Email sending flow completed successfully! ===")
}

// createEmailNotificationRequest creates a sample email notification request
func createEmailNotificationRequest() *models.EmailNotificationRequest {
	// Create content with subject and email body
	content := map[string]interface{}{
		"subject":    "Welcome to Our Service!",
		"email_body": "<h1>Welcome!</h1><p>Thank you for joining our service. We're excited to have you on board!</p><p>Best regards,<br>The Team</p>",
	}

	// Create email notification request
	notification := &models.EmailNotificationRequest{
		ID:         "email-123",
		Type:       "email",
		Content:    content,
		Recipients: []string{"user@example.com"},
		From: &models.EmailSender{
			Email: "noreply@example.com",
		},
	}

	fmt.Printf("Created EmailNotificationRequest:\n")
	fmt.Printf("  ID: %s\n", notification.ID)
	fmt.Printf("  Type: %s\n", notification.Type)
	fmt.Printf("  Recipients: %v\n", notification.Recipients)
	fmt.Printf("  From: %s\n", notification.From.Email)
	fmt.Printf("  Subject: %s\n", content["subject"])
	fmt.Printf("  Body: %s\n", content["email_body"])
	fmt.Println()

	return notification
}

// convertToJSONPayload converts the email notification to the JSON format
// that would be sent through Kafka channels
func convertToJSONPayload(notification *models.EmailNotificationRequest) string {
	// Create the notification data structure that matches what the handler creates
	notificationData := map[string]interface{}{
		"notification_id": notification.ID,
		"type":            notification.Type,
		"content":         notification.Content,
		"recipient": map[string]interface{}{
			"user_id":   "user-123",
			"email":     notification.Recipients[0],
			"full_name": "John Doe",
		},
		"from": map[string]interface{}{
			"email": notification.From.Email,
		},
		"created_at": time.Now(),
	}

	// Convert to JSON
	payload, err := json.Marshal(notificationData)
	if err != nil {
		log.Fatalf("Failed to marshal notification data: %v", err)
	}

	fmt.Printf("JSON Payload:\n%s\n", string(payload))
	fmt.Println()

	return string(payload)
}

// createNotificationMessage creates a notification message from the JSON payload
func createNotificationMessage(payload string) consumers.NotificationMessage {
	message := consumers.NotificationMessage{
		Type:      consumers.EmailNotification,
		Payload:   payload,
		ID:        "email-123",
		Timestamp: time.Now().Unix(),
	}

	fmt.Printf("Created NotificationMessage:\n")
	fmt.Printf("  Type: %s\n", message.Type)
	fmt.Printf("  ID: %s\n", message.ID)
	fmt.Printf("  Timestamp: %d\n", message.Timestamp)
	fmt.Println()

	return message
}

// processEmailNotification processes the email notification using the email processor
func processEmailNotification(message consumers.NotificationMessage) {
	// Create email processor
	processor := consumers.NewEmailProcessor()

	// Process the notification
	ctx := context.Background()
	err := processor.ProcessNotification(ctx, message)

	if err != nil {
		fmt.Printf("Error processing email notification: %v\n", err)
		return
	}

	fmt.Println("Email notification processed successfully!")
	fmt.Println("Note: In a real environment with SMTP configuration, this would send an actual email.")
	fmt.Println("In this example, it uses the mock email service since no SMTP credentials are configured.")
}

// Example of how to use the email processor with a custom email service
func exampleWithCustomEmailService() {
	fmt.Println("\n=== Example with Custom Email Service ===")

	// Create a custom email service (could be SendGrid, AWS SES, etc.)
	customService := &customEmailService{}

	// Create email processor with custom service
	processor := consumers.NewEmailProcessorWithService(customService)

	// Create test message
	message := consumers.NotificationMessage{
		Type:      consumers.EmailNotification,
		Payload:   `{"notification_id":"test-456","type":"email","content":{"subject":"Test","email_body":"Test body"},"recipient":{"email":"test@example.com"},"created_at":"2023-01-01T00:00:00Z"}`,
		ID:        "test-456",
		Timestamp: time.Now().Unix(),
	}

	// Process notification
	ctx := context.Background()
	err := processor.ProcessNotification(ctx, message)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Custom email service processed notification successfully!")
	}
}

// customEmailService is an example of a custom email service implementation
type customEmailService struct{}

func (c *customEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	// This would integrate with your preferred email service
	// Examples: SendGrid, AWS SES, Mailgun, etc.

	fmt.Println("Custom email service: Sending email...")

	// Return success response
	return &models.EmailResponse{
		ID:      "custom-123",
		Status:  "sent",
		Message: "Email sent via custom service",
		SentAt:  time.Now(),
		Channel: "email",
	}, nil
}
