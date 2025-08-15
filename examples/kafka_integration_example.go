package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gaurav2721/notification-service/services"
)

func RunKafkaIntegrationExample() {
	// Create service factory
	factory := services.NewServiceFactory()

	// Create service container (consumer manager is started automatically)
	container := services.NewServiceContainer()

	// Alternative: Create services manually using factory
	fmt.Println("🔧 Creating services using factory...")

	// Create Kafka service using factory
	kafkaService, err := factory.NewKafkaService()
	if err != nil {
		log.Fatalf("Failed to create Kafka service: %v", err)
	}

	// Create consumer manager using factory
	consumerManager := factory.NewConsumerManager(kafkaService)

	// Create notification manager using factory (with only Kafka service)
	notificationService := factory.NewNotificationManagerWithKafkaOnly(kafkaService)

	// Consumer manager is already started in the container
	fmt.Println("✅ Consumer manager started automatically")
	fmt.Println("✅ Kafka service initialized")
	fmt.Println("✅ All worker pools are running")
	fmt.Println("✅ Factory-created services ready")

	// Create context for operations
	ctx := context.Background()
	defer container.Shutdown(ctx)

	// Get the notification service from container (or use factory-created one)
	containerNotificationService := container.GetNotificationService()

	// Example: Send an email notification using container service
	emailNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		ScheduledAt *time.Time
	}{
		ID:   "email-001",
		Type: "email",
		Content: map[string]interface{}{
			"subject": "Welcome to our service!",
			"body":    "Thank you for signing up.",
		},
		Recipients: []string{"user@example.com"},
	}

	// Send the notification using container service
	response, err := containerNotificationService.SendNotification(ctx, emailNotification)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	} else {
		fmt.Printf("📧 Email notification sent (container): %+v\n", response)
	}

	// Example: Send a Slack notification using factory-created service
	slackNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		ScheduledAt *time.Time
	}{
		ID:   "slack-001",
		Type: "slack",
		Content: map[string]interface{}{
			"channel": "#general",
			"message": "Hello from the notification service!",
		},
		Recipients: []string{"#general"},
	}

	// Send the notification using factory-created service
	response, err = notificationService.SendNotification(ctx, slackNotification)
	if err != nil {
		log.Printf("Failed to send Slack notification: %v", err)
	} else {
		fmt.Printf("💬 Slack notification sent (factory): %+v\n", response)
	}

	// Example: Send an iOS push notification
	iosNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		ScheduledAt *time.Time
	}{
		ID:   "ios-push-001",
		Type: "ios_push",
		Content: map[string]interface{}{
			"title": "New Message",
			"body":  "You have a new message",
			"badge": 1,
			"sound": "default",
		},
		Recipients: []string{"device-token-123"},
	}

	// Send the notification
	response, err = notificationService.SendNotification(ctx, iosNotification)
	if err != nil {
		log.Printf("Failed to send iOS push notification: %v", err)
	} else {
		fmt.Printf("📱 iOS push notification sent: %+v\n", response)
	}

	// Example: Send an Android push notification
	androidNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		ScheduledAt *time.Time
	}{
		ID:   "android-push-001",
		Type: "android_push",
		Content: map[string]interface{}{
			"title":    "New Message",
			"body":     "You have a new message",
			"priority": "high",
		},
		Recipients: []string{"device-token-456"},
	}

	// Send the notification
	response, err = notificationService.SendNotification(ctx, androidNotification)
	if err != nil {
		log.Printf("Failed to send Android push notification: %v", err)
	} else {
		fmt.Printf("🤖 Android push notification sent: %+v\n", response)
	}

	// Wait a bit to see the consumers process the notifications
	fmt.Println("\n⏳ Waiting for consumers to process notifications...")
	time.Sleep(5 * time.Second)

	// Get consumer status from container
	containerConsumerManager := container.GetConsumerManager()
	status := containerConsumerManager.GetStatus()
	fmt.Println("\n📊 Consumer Status (Container):")
	for notificationType, isRunning := range status {
		fmt.Printf("  %s: %t\n", notificationType, isRunning)
	}

	// Demonstrate factory-created consumer manager
	fmt.Println("\n🔧 Factory-created Consumer Manager:")
	fmt.Printf("  Email Workers: %d\n", consumerManager.GetWorkerPool(services.EmailNotification).GetWorkerCount())
	fmt.Printf("  Slack Workers: %d\n", consumerManager.GetWorkerPool(services.SlackNotification).GetWorkerCount())

	fmt.Println("\n✅ Example completed successfully!")
}

func main() {
	RunKafkaIntegrationExample()
}
