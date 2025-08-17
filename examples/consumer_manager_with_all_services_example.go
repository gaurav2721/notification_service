//go:build example
// +build example

package main

import (
	"context"
	"fmt"
)

// Example demonstrating the consumer manager with all service dependencies
func RunConsumerManagerWithAllServicesExample() {
	fmt.Println("🚀 Running Consumer Manager with All Services Example")

	// This example shows how to create a consumer manager with all service dependencies
	// including email, slack, APNS, and FCM services

	fmt.Println("✅ Consumer manager refactored to support all service dependencies")
	fmt.Println("📧 Email service: Supports SMTP and custom email providers")
	fmt.Println("💬 Slack service: Supports Slack API and webhooks")
	fmt.Println("🍎 APNS service: Supports Apple Push Notification Service")
	fmt.Println("🤖 FCM service: Supports Firebase Cloud Messaging")
	fmt.Println("📊 Kafka service: Supports message queuing and distribution")

	fmt.Println("\n🔧 Key Features:")
	fmt.Println("  - Dependency injection for all services")
	fmt.Println("  - Service reuse across the application")
	fmt.Println("  - Better testability with mock services")
	fmt.Println("  - Flexible service configuration")
	fmt.Println("  - Backward compatibility maintained")

	fmt.Println("\n📝 Usage Examples:")
	fmt.Println("  1. Basic usage with default services")
	fmt.Println("  2. Advanced usage with custom services")
	fmt.Println("  3. Testing with mock services")
	fmt.Println("  4. Integration with service factory")
}

// Example of creating processors with custom services
func ExampleCustomProcessorsWithAllServices() {
	fmt.Println("\n🔧 Example: Custom Processors with All Services")

	fmt.Println("✅ Custom services available:")
	fmt.Println("  - Custom Email Service")
	fmt.Println("  - Custom Slack Service")
	fmt.Println("  - Custom APNS Service")
	fmt.Println("  - Custom FCM Service")

	// Example processor creation with custom services
	fmt.Println("\n📱 Processor Creation Examples:")
	fmt.Println("  - NewEmailProcessorWithService(customEmailService)")
	fmt.Println("  - NewSlackProcessorWithService(customSlackService)")
	fmt.Println("  - NewIOSPushProcessorWithService(customAPNSService)")
	fmt.Println("  - NewAndroidPushProcessorWithService(customFCMService)")

	// Example consumer manager creation
	fmt.Println("\n🏗️  Consumer Manager Creation:")
	fmt.Println("  - NewConsumerManagerWithServices(email, slack, apns, fcm, kafka, config)")
	fmt.Println("  - All services injected as dependencies")
	fmt.Println("  - Processors automatically use injected services")
}

// Example of testing with mock services
func ExampleTestingWithMockServices() {
	fmt.Println("\n🧪 Example: Testing with Mock Services")

	fmt.Println("✅ Mock services available for testing:")
	fmt.Println("  - MockEmailService")
	fmt.Println("  - MockSlackService")
	fmt.Println("  - MockAPNSService")
	fmt.Println("  - MockFCMService")

	fmt.Println("\n📋 Testing Benefits:")
	fmt.Println("  - Isolated unit tests")
	fmt.Println("  - No external dependencies")
	fmt.Println("  - Fast test execution")
	fmt.Println("  - Predictable test results")
	fmt.Println("  - Easy to verify service calls")
}

// Example of service factory usage
func ExampleServiceFactoryUsage() {
	fmt.Println("\n🏭 Example: Service Factory Usage")

	fmt.Println("✅ Factory methods available:")
	fmt.Println("  - NewEmailService()")
	fmt.Println("  - NewSlackService()")
	fmt.Println("  - NewAPNSService()")
	fmt.Println("  - NewFCMService()")
	fmt.Println("  - NewKafkaService()")
	fmt.Println("  - NewConsumerManagerWithServices(...)")

	fmt.Println("\n🔧 Factory Benefits:")
	fmt.Println("  - Consistent service creation")
	fmt.Println("  - Centralized configuration")
	fmt.Println("  - Easy dependency management")
	fmt.Println("  - Simplified testing setup")
}

// Custom service implementations for demonstration
type CustomEmailService struct{}

func (c *CustomEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	fmt.Println("📧 Custom email service: Sending email...")
	return map[string]interface{}{
		"status":  "sent",
		"service": "custom_email",
	}, nil
}

type CustomSlackService struct{}

func (c *CustomSlackService) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	fmt.Println("💬 Custom slack service: Sending slack message...")
	return map[string]interface{}{
		"status":  "sent",
		"service": "custom_slack",
	}, nil
}

type CustomAPNSService struct{}

func (c *CustomAPNSService) SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	fmt.Println("🍎 Custom APNS service: Sending iOS push notification...")
	return map[string]interface{}{
		"status":  "sent",
		"service": "custom_apns",
	}, nil
}

type CustomFCMService struct{}

func (c *CustomFCMService) SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	fmt.Println("🤖 Custom FCM service: Sending Android push notification...")
	return map[string]interface{}{
		"status":  "sent",
		"service": "custom_fcm",
	}, nil
}

// Example of notification payloads
func ExampleNotificationPayloads() {
	fmt.Println("\n📦 Example: Notification Payloads")

	fmt.Println("📧 Email Notification Payload:")
	fmt.Println(`{
  "notification_id": "email-123",
  "type": "email",
  "content": {
    "subject": "Welcome!",
    "email_body": "<h1>Welcome to our service!</h1>"
  },
  "recipient": {
    "email": "user@example.com"
  },
  "from": {
    "email": "noreply@example.com"
  }
}`)

	fmt.Println("\n💬 Slack Notification Payload:")
	fmt.Println(`{
  "notification_id": "slack-456",
  "type": "slack",
  "content": {
    "text": "Hello from our service!"
  },
  "channel": "#general"
}`)

	fmt.Println("\n🍎 iOS Push Notification Payload:")
	fmt.Println(`{
  "notification_id": "ios-789",
  "type": "ios_push",
  "content": {
    "title": "New Message",
    "body": "You have a new message"
  },
  "recipients": ["ios_device_token_123", "ios_device_token_456"]
}`)

	fmt.Println("\n🤖 Android Push Notification Payload:")
	fmt.Println(`{
  "notification_id": "android-101",
  "type": "android_push",
  "content": {
    "title": "New Message",
    "body": "You have a new message"
  },
  "recipients": ["android_device_token_123", "android_device_token_456"]
}`)
}

// Example of migration guide
func ExampleMigrationGuide() {
	fmt.Println("\n🔄 Migration Guide")

	fmt.Println("✅ For Existing Code:")
	fmt.Println("  - No changes required")
	fmt.Println("  - Existing NewConsumerManager() continues to work")
	fmt.Println("  - Gradual migration possible")

	fmt.Println("\n🆕 For New Code:")
	fmt.Println("  - Use NewConsumerManagerWithServices()")
	fmt.Println("  - Inject all required services")
	fmt.Println("  - Leverage service factory")
	fmt.Println("  - Use mock services for testing")

	fmt.Println("\n🧪 For Testing:")
	fmt.Println("  - Create mock services")
	fmt.Println("  - Inject mocks into consumer manager")
	fmt.Println("  - Test processors in isolation")
	fmt.Println("  - Verify service interactions")
}
