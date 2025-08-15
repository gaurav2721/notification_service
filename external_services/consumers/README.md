# Consumer Worker Pool System

This package provides a robust consumer worker pool system for processing different types of notifications (email, slack, iOS push notifications, and Android push notifications) from Kafka channels.

## Overview

The consumer system consists of:

1. **ConsumerWorker**: Individual workers that process notifications from channels
2. **ConsumerWorkerPool**: Pools of workers for specific notification types
3. **ConsumerManager**: Orchestrates all worker pools
4. **NotificationProcessor**: Interface for processing different notification types

## Architecture

```
ConsumerManager
├── EmailWorkerPool (configurable number of workers)
│   ├── Worker 1
│   ├── Worker 2
│   └── Worker N
├── SlackWorkerPool (configurable number of workers)
│   ├── Worker 1
│   ├── Worker 2
│   └── Worker N
├── IOSPushWorkerPool (configurable number of workers)
│   ├── Worker 1
│   ├── Worker 2
│   └── Worker N
└── AndroidPushWorkerPool (configurable number of workers)
    ├── Worker 1
    ├── Worker 2
    └── Worker N
```

## Configuration

Worker counts are configured via environment variables:

```bash
# Consumer Worker Pool Configuration
EMAIL_WORKER_COUNT=5
SLACK_WORKER_COUNT=3
IOS_PUSH_WORKER_COUNT=3
ANDROID_PUSH_WORKER_COUNT=3
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/gaurav2721/notification-service/external_services/consumers"
    "github.com/gaurav2721/notification-service/external_services/kafka"
)

func main() {
    // Initialize your Kafka service
    var kafkaService kafka.KafkaService
    // kafkaService = your_kafka_implementation.NewKafkaService()

    // Create consumer manager from environment configuration
    consumerManager := consumers.NewConsumerManagerFromEnv(kafkaService)

    // Create a context that can be cancelled
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Initialize the consumer manager
    if err := consumerManager.Initialize(ctx); err != nil {
        log.Fatalf("Failed to initialize consumer manager: %v", err)
    }

    // Start all worker pools
    if err := consumerManager.Start(ctx); err != nil {
        log.Fatalf("Failed to start consumer manager: %v", err)
    }

    // Your application logic here...

    // Graceful shutdown
    if err := consumerManager.Stop(); err != nil {
        log.Printf("Error stopping consumer manager: %v", err)
    }
}
```

### Custom Configuration

```go
config := consumers.ConsumerConfig{
    EmailWorkerCount:       10, // More workers for email
    SlackWorkerCount:       2,  // Fewer workers for slack
    IOSPushWorkerCount:     5,  // Medium workers for iOS
    AndroidPushWorkerCount: 5,  // Medium workers for Android
    KafkaService:           kafkaService,
}

consumerManager := consumers.NewConsumerManager(config)
```

### Monitoring Worker Pools

```go
// Get status of all worker pools
status := consumerManager.GetStatus()
for notificationType, isRunning := range status {
    log.Printf("Worker pool for %s: %v", notificationType, isRunning)
}

// Get a specific worker pool
if emailPool, err := consumerManager.GetWorkerPool(consumers.EmailNotification); err == nil {
    log.Printf("Email worker pool has %d workers", emailPool.GetWorkerCount())
}
```

## Notification Types

The system supports four notification types:

- `EmailNotification`: Email notifications
- `SlackNotification`: Slack messages
- `IOSPushNotification`: iOS push notifications
- `AndroidPushNotification`: Android push notifications

## Extending the System

### Adding New Notification Types

1. Add a new constant to `NotificationType`:
```go
const (
    // ... existing types ...
    NewNotificationType NotificationType = "new_type"
)
```

2. Create a new processor:
```go
type newNotificationProcessor struct{}

func NewNewNotificationProcessor() NotificationProcessor {
    return &newNotificationProcessor{}
}

func (np *newNotificationProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
    // Implement your processing logic
    return nil
}

func (np *newNotificationProcessor) GetNotificationType() NotificationType {
    return NewNotificationType
}
```

3. Add the processor to the consumer manager's initialization methods.

### Custom Processors

Each notification type has its own processor that implements the `NotificationProcessor` interface:

```go
type NotificationProcessor interface {
    ProcessNotification(ctx context.Context, message NotificationMessage) error
    GetNotificationType() NotificationType
}
```

## Error Handling

The system includes robust error handling:

- Individual worker failures don't stop the entire pool
- Graceful shutdown with context cancellation
- Thread-safe operations with mutex protection
- Comprehensive logging for debugging

## Performance Considerations

- Worker pools are isolated by notification type
- Configurable worker counts allow for load balancing
- Non-blocking channel operations
- Efficient context-based cancellation

## Testing

The system is designed to be easily testable:

- All interfaces are mockable
- Dependency injection through configuration
- Isolated worker pools for unit testing
- Context-based cancellation for integration tests 