package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gaurav2721/notification-service/services"
)

func RunKafkaIntegrationExample() {
	// Create service container
	container := services.NewServiceContainer()

	// Start the consumer manager
	ctx := context.Background()
	if err := container.StartConsumerManager(ctx); err != nil {
		log.Fatalf("Failed to start consumer manager: %v", err)
	}
	defer container.Shutdown(ctx)

	fmt.Println("✅ Consumer manager started successfully")
	fmt.Println("✅ Kafka service initialized")
	fmt.Println("✅ All worker pools are running")

	// Get the notification service
	notificationService := container.GetNotificationService()

	// Example: Send an email notification
	emailNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	}{
		ID:   "email-001",
		Type: "email",
		Content: map[string]interface{}{
			"subject": "Welcome to our service!",
			"body":    "Thank you for signing up.",
		},
		Recipients: []string{"user@example.com"},
		Metadata: map[string]interface{}{
			"priority": "high",
		},
	}

	// Send the notification
	response, err := notificationService.SendNotification(ctx, emailNotification)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	} else {
		fmt.Printf("📧 Email notification sent: %+v\n", response)
	}

	// Example: Send a Slack notification
	slackNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	}{
		ID:   "slack-001",
		Type: "slack",
		Content: map[string]interface{}{
			"channel": "#general",
			"message": "Hello from the notification service!",
		},
		Recipients: []string{"#general"},
		Metadata: map[string]interface{}{
			"priority": "normal",
		},
	}

	// Send the notification
	response, err = notificationService.SendNotification(ctx, slackNotification)
	if err != nil {
		log.Printf("Failed to send Slack notification: %v", err)
	} else {
		fmt.Printf("💬 Slack notification sent: %+v\n", response)
	}

	// Example: Send an iOS push notification
	iosNotification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    interface{}
		Recipients  []string
		Metadata    map[string]interface{}
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
		Metadata: map[string]interface{}{
			"priority": "high",
		},
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
		Metadata    map[string]interface{}
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
		Metadata: map[string]interface{}{
			"priority": "high",
		},
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

	// Get consumer status
	consumerManager := container.GetConsumerManager()
	status := consumerManager.GetStatus()
	fmt.Println("\n📊 Consumer Status:")
	for notificationType, isRunning := range status {
		fmt.Printf("  %s: %t\n", notificationType, isRunning)
	}

	fmt.Println("\n✅ Example completed successfully!")
}

func main() {
	RunKafkaIntegrationExample()
}
