//go:build example
// +build example

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gaurav2721/notification-service/external_services/consumers"
	"github.com/gaurav2721/notification-service/external_services/email"
	"github.com/gaurav2721/notification-service/external_services/kafka"
	"github.com/gaurav2721/notification-service/external_services/slack"
)

func RunConsumerManagerWithServicesExample() {
	fmt.Println("🚀 Running Consumer Manager with Services Example")

	// Create services
	emailService := email.NewEmailService()
	slackService := slack.NewSlackService()
	kafkaService, err := kafka.NewKafkaService()
	if err != nil {
		log.Fatalf("Failed to create Kafka service: %v", err)
	}

	// Create consumer configuration
	config := consumers.ConsumerConfig{
		EmailWorkerCount:       3,
		SlackWorkerCount:       2,
		IOSPushWorkerCount:     2,
		AndroidPushWorkerCount: 2,
	}

	// Create consumer manager with service dependencies
	consumerManager := consumers.NewConsumerManagerWithServices(
		emailService,
		slackService,
		kafkaService,
		config,
	)

	// Initialize and start the consumer manager
	ctx := context.Background()
	if err := consumerManager.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize consumer manager: %v", err)
	}

	if err := consumerManager.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer manager: %v", err)
	}

	fmt.Println("✅ Consumer manager started with service dependencies")

	// Get status of all worker pools
	status := consumerManager.GetStatus()
	fmt.Println("📊 Worker pool status:")
	for notificationType, isRunning := range status {
		fmt.Printf("  - %s: %t\n", notificationType, isRunning)
	}

	// Example: Send a test notification to email channel
	emailChannel := kafkaService.GetEmailChannel()
	testMessage := `{
		"notification_id": "test-123",
		"type": "email",
		"content": {
			"subject": "Test Email",
			"email_body": "<h1>Hello from refactored consumer manager!</h1>"
		},
		"recipient": {
			"email": "test@example.com"
		},
		"from": {
			"email": "noreply@example.com"
		}
	}`

	select {
	case emailChannel <- testMessage:
		fmt.Println("✅ Test email notification sent to Kafka channel")
	case <-time.After(5 * time.Second):
		fmt.Println("⚠️  Timeout sending test notification")
	}

	// Example: Send a test notification to slack channel
	slackChannel := kafkaService.GetSlackChannel()
	slackMessage := `{
		"notification_id": "test-456",
		"type": "slack",
		"content": {
			"text": "Hello from refactored consumer manager!"
		},
		"channel": "#general"
	}`

	select {
	case slackChannel <- slackMessage:
		fmt.Println("✅ Test slack notification sent to Kafka channel")
	case <-time.After(5 * time.Second):
		fmt.Println("⚠️  Timeout sending test slack notification")
	}

	// Wait a bit for processing
	time.Sleep(2 * time.Second)

	// Stop the consumer manager
	if err := consumerManager.Stop(); err != nil {
		log.Printf("Error stopping consumer manager: %v", err)
	} else {
		fmt.Println("✅ Consumer manager stopped gracefully")
	}

	// Close Kafka service
	kafkaService.Close()
	fmt.Println("✅ Kafka service closed")
}

// Example of creating processors with custom services
func ExampleCustomProcessors() {
	fmt.Println("\n🔧 Example: Custom Processors with Services")

	// Create custom email service
	customEmailService := &CustomEmailService{}
	customSlackService := &CustomSlackService{}

	// Create processors with custom services
	emailProcessor := consumers.NewEmailProcessorWithService(customEmailService)
	slackProcessor := consumers.NewSlackProcessorWithService(customSlackService)

	fmt.Printf("Email processor type: %s\n", emailProcessor.GetNotificationType())
	fmt.Printf("Slack processor type: %s\n", slackProcessor.GetNotificationType())

	// Test processing a notification
	ctx := context.Background()
	testMessage := consumers.NotificationMessage{
		Type:      consumers.EmailNotification,
		Payload:   `{"content":{"subject":"Test","email_body":"Test body"},"recipient":{"email":"test@example.com"}}`,
		ID:        "test-789",
		Timestamp: time.Now().Unix(),
	}

	if err := emailProcessor.ProcessNotification(ctx, testMessage); err != nil {
		fmt.Printf("Error processing notification: %v\n", err)
	} else {
		fmt.Println("✅ Custom email processor processed notification successfully")
	}
}

// Custom service implementations for demonstration
type CustomEmailService struct{}

func (c *CustomEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	fmt.Println("📧 Custom email service: Sending email...")
	return map[string]interface{}{
		"status":  "sent",
		"service": "custom",
	}, nil
}

type CustomSlackService struct{}

func (c *CustomSlackService) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	fmt.Println("💬 Custom slack service: Sending slack message...")
	return map[string]interface{}{
		"status":  "sent",
		"service": "custom",
	}, nil
}
