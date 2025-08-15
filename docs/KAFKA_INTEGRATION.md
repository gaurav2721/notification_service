# Kafka Integration Implementation

This document describes the implementation of Kafka integration for the notification service, which enables asynchronous processing of notifications through Kafka channels and consumer workers.

## Overview

The implementation consists of three main tasks:

1. **Initialize Kafka service from external_services**
2. **Initialize the consumers in external_services**
3. **Pass the Kafka channels in notification manager**

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Notification  │    │   Kafka Service  │    │   Consumer      │
│   Manager       │───▶│   (Channels)     │───▶│   Manager       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │   Email Channel  │    │   Email Workers │
                       │   Slack Channel  │    │   Slack Workers │
                       │   iOS Channel    │    │   iOS Workers   │
                       │   Android Channel│    │   Android Workers│
                       └──────────────────┘    └─────────────────┘
```

## Implementation Details

### 1. Kafka Service Integration

**File**: `services/service_provider.go`

- Added `kafkaService` field to `ServiceContainer`
- Initialized Kafka service using `kafka.NewKafkaService()`
- Added `GetKafkaService()` method to access the Kafka service
- Updated `Shutdown()` method to properly close Kafka service

**Key Changes**:
```go
type ServiceContainer struct {
    // ... existing fields ...
    kafkaService        kafka.KafkaService
    consumerManager     consumers.ConsumerManager
    // ... existing fields ...
}
```

### 2. Consumer Manager Integration

**File**: `services/service_provider.go`

- Added `consumerManager` field to `ServiceContainer`
- Initialized consumer manager with Kafka service configuration
- **Consumer manager is now started automatically during initialization**
- Added `GetConsumerManager()` method to access the consumer manager
- Removed `StartConsumerManager()` method since it's no longer needed

**Key Changes**:
```go
// Initialize consumer manager using factory
c.consumerManager = factory.NewConsumerManager(c.kafkaService)

// Start the consumer manager immediately
ctx := context.Background()
if err := c.consumerManager.Initialize(ctx); err != nil {
    panic("Failed to initialize consumer manager: " + err.Error())
}
if err := c.consumerManager.Start(ctx); err != nil {
    panic("Failed to start consumer manager: " + err.Error())
}
```

### 3. Notification Manager Updates

**File**: `notification_manager/notification_manager.go`

- Added `kafkaService` field to `NotificationManagerImpl`
- Updated constructor to accept Kafka service parameter
- **Modified to only receive Kafka service since it only pushes to channels**
- Modified `SendNotification()` method to send notifications to Kafka channels

**Key Changes**:
```go
// Send notification to Kafka channel based on type
switch notif.Type {
case "email":
    if kafkaService, ok := nm.kafkaService.(interface {
        GetEmailChannel() chan string
    }); ok {
        message := fmt.Sprintf("Email notification: %s", notif.ID)
        select {
        case kafkaService.GetEmailChannel() <- message:
            // Message sent successfully
        default:
            // Channel is full, handle accordingly
        }
    }
    // ... return response
}
```

### 4. Service Factory Updates

**File**: `services/service_factory.go`

- Added Kafka service import and type alias
- Added ConsumerManager import and type alias
- Updated `NewNotificationManager()` method to accept Kafka service parameter
- Added `NewKafkaService()` method for creating Kafka service instances
- Added `NewConsumerManager()` method for creating consumer manager instances
- Added `NewConsumerManagerWithConfig()` method for custom configuration
- Added `NewConsumerManagerFromEnv()` method for environment-based configuration
- Added multiple notification manager factory methods for different dependency combinations

**Key Changes**:
```go
type KafkaService = kafka.KafkaService
type ConsumerManager = consumers.ConsumerManager

func (f *ServiceFactory) NewKafkaService() (KafkaService, error) {
    return kafka.NewKafkaService()
}

func (f *ServiceFactory) NewConsumerManager(kafkaService KafkaService) ConsumerManager {
    config := consumers.ConsumerConfig{
        KafkaService: kafkaService,
    }
    return consumers.NewConsumerManager(config)
}

func (f *ServiceFactory) NewNotificationManager(
    emailService EmailService,
    slackService SlackService,
    inAppService InAppService,
    kafkaService KafkaService,
) NotificationManager {
    return notification_manager.NewNotificationManagerWithDefaultTemplate(
        emailService, slackService, inAppService, nil, kafkaService, nil)
}
```

**Factory Methods Available**:
- `NewKafkaService()` - Creates a new Kafka service instance
- `NewConsumerManager(kafkaService)` - Creates a consumer manager with default config
- `NewConsumerManagerWithConfig(config)` - Creates a consumer manager with custom config
- `NewConsumerManagerFromEnv(kafkaService)` - Creates a consumer manager from environment variables
- `NewNotificationManager(...)` - Creates notification manager with basic dependencies
- `NewNotificationManagerWithKafkaOnly(kafkaService)` - Creates notification manager with only Kafka service (for channel pushing)
- `NewNotificationManagerWithUserService(...)` - Creates notification manager with user service
- `NewNotificationManagerWithScheduler(...)` - Creates notification manager with scheduler
- `NewNotificationManagerComplete(...)` - Creates notification manager with all dependencies

### 5. Consumer Processor Implementation

**File**: `external_services/consumers/manager.go`

- Added processor factory functions (`NewEmailProcessor`, `NewSlackProcessor`, etc.)
- Added processor implementations for each notification type
- Each processor implements the `NotificationProcessor` interface

**Key Changes**:
```go
// Processor factory functions
func NewEmailProcessor() NotificationProcessor {
    return &emailProcessor{}
}

// Processor implementations
type emailProcessor struct{}

func (ep *emailProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
    // TODO: Implement actual email sending logic
    return nil
}
```

## Usage Example

**File**: `examples/kafka_integration_example.go`

The example demonstrates how to:

1. Create a service container
2. Create services using factory methods
3. Start the consumer manager
4. Send notifications of different types
5. Monitor consumer status

```go
func RunKafkaIntegrationExample() {
    // Create service factory
    factory := services.NewServiceFactory()
    
    // Create service container (consumer manager starts automatically)
    container := services.NewServiceContainer()

    // Alternative: Create services manually using factory
    kafkaService, err := factory.NewKafkaService()
    if err != nil {
        log.Fatalf("Failed to create Kafka service: %v", err)
    }
    
    consumerManager := factory.NewConsumerManager(kafkaService)
    notificationService := factory.NewNotificationManagerWithKafkaOnly(kafkaService)

    // Consumer manager is already started in the container
    ctx := context.Background()
    defer container.Shutdown(ctx)

    // Send notifications
    response, err := notificationService.SendNotification(ctx, emailNotification)
    // ... handle response
}
```

## Supported Notification Types

The implementation supports the following notification types:

1. **email** - Sent to email channel
2. **slack** - Sent to Slack channel  
3. **ios_push** - Sent to iOS push notification channel
4. **android_push** - Sent to Android push notification channel
5. **in_app** - Direct processing (not through Kafka)

## Configuration

The Kafka service and consumer manager can be configured through environment variables:

- `EMAIL_CHANNEL_BUFFER_SIZE` - Buffer size for email channel (default: 100)
- `SLACK_CHANNEL_BUFFER_SIZE` - Buffer size for Slack channel (default: 100)
- `IOS_PUSH_CHANNEL_BUFFER_SIZE` - Buffer size for iOS push channel (default: 100)
- `ANDROID_PUSH_CHANNEL_BUFFER_SIZE` - Buffer size for Android push channel (default: 100)
- `EMAIL_WORKER_COUNT` - Number of email workers (default: 5)
- `SLACK_WORKER_COUNT` - Number of Slack workers (default: 3)
- `IOS_PUSH_WORKER_COUNT` - Number of iOS push workers (default: 3)
- `ANDROID_PUSH_WORKER_COUNT` - Number of Android push workers (default: 3)

## Benefits

1. **Asynchronous Processing**: Notifications are queued and processed asynchronously
2. **Scalability**: Multiple workers can process notifications in parallel
3. **Reliability**: Kafka channels provide buffering and fault tolerance
4. **Flexibility**: Easy to add new notification types and processors
5. **Monitoring**: Consumer status can be monitored and managed
6. **Factory Pattern**: Centralized service creation with multiple configuration options

## Future Enhancements

1. **JSON Serialization**: Implement proper JSON marshaling for notification messages
2. **Error Handling**: Add comprehensive error handling for channel full scenarios
3. **Metrics**: Add metrics and monitoring for Kafka channels and consumers
4. **Retry Logic**: Implement retry mechanisms for failed notifications
5. **Dead Letter Queue**: Add dead letter queue for failed notifications
6. **Configuration**: Add more configuration options for Kafka service
7. **Logging**: Add structured logging for better observability

## Testing

To test the implementation:

```bash
# Build the project
go build .

# Run the example
go run examples/kafka_integration_example.go
```

The example will demonstrate the complete flow from sending notifications to processing them through Kafka channels and consumer workers.