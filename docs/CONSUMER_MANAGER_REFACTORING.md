# Consumer Manager Refactoring: Service Dependencies

This document explains the refactoring changes made to the consumer manager to accept email and slack services as dependencies, allowing processors to utilize the same services for sending notifications.

## Overview

The consumer manager has been refactored to support dependency injection of email and slack services, enabling better testability, flexibility, and service reuse across the application.

## Key Changes

### 1. Updated ConsumerConfig Interface

**File**: `external_services/consumers/interfaces.go`

Added service dependencies to the `ConsumerConfig` struct:

```go
type ConsumerConfig struct {
    EmailWorkerCount       int `json:"email_worker_count" env:"EMAIL_WORKER_COUNT" env-default:"5"`
    SlackWorkerCount       int `json:"slack_worker_count" env:"SLACK_WORKER_COUNT" env-default:"3"`
    IOSPushWorkerCount     int `json:"ios_push_worker_count" env:"IOS_PUSH_WORKER_COUNT" env-default:"3"`
    AndroidPushWorkerCount int `json:"android_push_worker_count" env:"ANDROID_PUSH_WORKER_COUNT" env-default:"3"`

    // Service dependencies
    EmailService interface {
        SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
    }
    SlackService interface {
        SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
    }

    // Kafka service interface for getting channels
    KafkaService interface {
        GetEmailChannel() chan string
        GetSlackChannel() chan string
        GetIOSPushNotificationChannel() chan string
        GetAndroidPushNotificationChannel() chan string
        Close()
    }
}
```

### 2. New Consumer Manager Constructor

**File**: `external_services/consumers/manager.go`

Added a new constructor that accepts service dependencies:

```go
// NewConsumerManagerWithServices creates a new consumer manager with service dependencies
func NewConsumerManagerWithServices(
    emailService interface {
        SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
    },
    slackService interface {
        SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
    },
    kafkaService interface {
        GetEmailChannel() chan string
        GetSlackChannel() chan string
        GetIOSPushNotificationChannel() chan string
        GetAndroidPushNotificationChannel() chan string
        Close()
    },
    config ConsumerConfig,
) ConsumerManager {
    // Update config with service dependencies
    config.EmailService = emailService
    config.SlackService = slackService
    config.KafkaService = kafkaService

    return &consumerManager{
        config:      config,
        workerPools: make(map[NotificationType]ConsumerWorkerPool),
        running:     false,
    }
}
```

### 3. Updated Processor Creation

**File**: `external_services/consumers/manager.go`

Modified the processor creation methods to use injected services:

```go
// createEmailWorkerPool creates the email worker pool
func (cm *consumerManager) createEmailWorkerPool() {
    var processor NotificationProcessor
    
    // Use injected email service if available, otherwise create default
    if cm.config.EmailService != nil {
        processor = NewEmailProcessorWithService(cm.config.EmailService)
    } else {
        processor = NewEmailProcessor()
    }
    
    pool := NewWorkerPool(
        EmailNotification,
        cm.config.KafkaService.GetEmailChannel(),
        processor,
        cm.config.EmailWorkerCount,
    )
    cm.workerPools[EmailNotification] = pool
}

// createSlackWorkerPool creates the slack worker pool
func (cm *consumerManager) createSlackWorkerPool() {
    var processor NotificationProcessor
    
    // Use injected slack service if available, otherwise create default
    if cm.config.SlackService != nil {
        processor = NewSlackProcessorWithService(cm.config.SlackService)
    } else {
        processor = NewSlackProcessor()
    }
    
    pool := NewWorkerPool(
        SlackNotification,
        cm.config.KafkaService.GetSlackChannel(),
        processor,
        cm.config.SlackWorkerCount,
    )
    cm.workerPools[SlackNotification] = pool
}
```

### 4. Enhanced Slack Processor

**File**: `external_services/consumers/slack_processor.go`

Updated the slack processor to accept service dependencies and implemented actual processing logic:

```go
// slackProcessor handles slack notification processing
type slackProcessor struct {
    slackService interface {
        SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
    }
}

// NewSlackProcessorWithService creates a new slack processor with a specific slack service
func NewSlackProcessorWithService(slackService interface {
    SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
    return &slackProcessor{
        slackService: slackService,
    }
}

// ProcessNotification processes a slack notification
func (sp *slackProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
    // If no slack service is available, just log and return
    if sp.slackService == nil {
        logrus.Warn("No slack service available, skipping slack notification")
        return nil
    }

    // Parse the payload to extract notification details
    var notificationData map[string]interface{}
    if err := json.Unmarshal([]byte(message.Payload), &notificationData); err != nil {
        return fmt.Errorf("failed to parse notification payload: %w", err)
    }

    // Extract content and create slack notification request
    content, ok := notificationData["content"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("invalid content data in notification payload")
    }

    // Create slack notification request
    slackNotification := &struct {
        ID         string                 `json:"id"`
        Type       string                 `json:"type"`
        Content    map[string]interface{} `json:"content"`
        Recipients []string               `json:"recipients"`
        Channel    string                 `json:"channel,omitempty"`
    }{
        ID:         message.ID,
        Type:       string(message.Type),
        Content:    content,
        Recipients: []string{},
        Channel:    channel,
    }

    // Send slack message using the slack service
    response, err := sp.slackService.SendSlackMessage(ctx, slackNotification)
    if err != nil {
        return fmt.Errorf("failed to send slack message: %w", err)
    }

    logrus.WithFields(logrus.Fields{
        "notification_id": message.ID,
        "response":        response,
    }).Info("Slack notification sent successfully")

    return nil
}
```

### 5. Updated Service Provider

**File**: `services/service_provider.go`

Modified the service provider to use the new constructor with service dependencies:

```go
// Initialize consumer manager using factory with environment configuration
logrus.Debug("Initializing consumer manager")
// Use the new constructor with service dependencies
config := consumers.ConsumerConfig{
    EmailWorkerCount:       getEnvAsInt("EMAIL_WORKER_COUNT", 5),
    SlackWorkerCount:       getEnvAsInt("SLACK_WORKER_COUNT", 3),
    IOSPushWorkerCount:     getEnvAsInt("IOS_PUSH_WORKER_COUNT", 3),
    AndroidPushWorkerCount: getEnvAsInt("ANDROID_PUSH_WORKER_COUNT", 3),
}
c.consumerManager = consumers.NewConsumerManagerWithServices(
    c.emailService,
    c.slackService,
    c.kafkaService,
    config,
)
```

### 6. Enhanced Service Factory

**File**: `services/service_factory.go`

Added a new factory method for creating consumer managers with service dependencies:

```go
// NewConsumerManagerWithServices creates a new consumer manager with service dependencies
func (f *ServiceFactory) NewConsumerManagerWithServices(
    emailService EmailService,
    slackService SlackService,
    kafkaService KafkaService,
    config ConsumerConfig,
) ConsumerManager {
    return consumers.NewConsumerManagerWithServices(emailService, slackService, kafkaService, config)
}
```

## Benefits of This Refactoring

### 1. **Dependency Injection**
- Services can now be injected into the consumer manager
- Better testability with mock services
- Flexible service configuration

### 2. **Service Reuse**
- The same email and slack services can be used across the application
- Consistent service behavior and configuration
- Reduced resource usage

### 3. **Better Testing**
- Easy to mock services for unit tests
- Isolated testing of processors
- Better integration test scenarios

### 4. **Flexibility**
- Can use different service implementations (e.g., SendGrid vs SMTP for email)
- Easy to switch between real and mock services
- Support for custom service implementations

### 5. **Backward Compatibility**
- Existing code continues to work with default constructors
- Gradual migration path available
- No breaking changes to existing interfaces

## Usage Examples

### Basic Usage with Default Services

```go
// Create consumer manager with default services
config := consumers.ConsumerConfig{
    EmailWorkerCount: 5,
    SlackWorkerCount: 3,
}
consumerManager := consumers.NewConsumerManager(config)
```

### Advanced Usage with Custom Services

```go
// Create custom services
emailService := &CustomEmailService{}
slackService := &CustomSlackService{}
kafkaService := &KafkaService{}

// Create consumer manager with custom services
config := consumers.ConsumerConfig{
    EmailWorkerCount: 5,
    SlackWorkerCount: 3,
}
consumerManager := consumers.NewConsumerManagerWithServices(
    emailService,
    slackService,
    kafkaService,
    config,
)
```

### Using Service Factory

```go
factory := services.NewServiceFactory()

// Create services
emailService := factory.NewEmailService()
slackService := factory.NewSlackService()
kafkaService, _ := factory.NewKafkaService()

// Create consumer manager with services
config := consumers.ConsumerConfig{
    EmailWorkerCount: 5,
    SlackWorkerCount: 3,
}
consumerManager := factory.NewConsumerManagerWithServices(
    emailService,
    slackService,
    kafkaService,
    config,
)
```

## Migration Guide

### For Existing Code

1. **No Changes Required**: Existing code using `NewConsumerManager()` continues to work
2. **Optional Enhancement**: Update to use `NewConsumerManagerWithServices()` for better dependency management
3. **Service Provider**: The service provider automatically uses the new constructor

### For New Code

1. **Use New Constructor**: Prefer `NewConsumerManagerWithServices()` for new implementations
2. **Inject Services**: Pass email and slack services as dependencies
3. **Use Factory**: Leverage the service factory for consistent service creation

## Testing

### Unit Testing with Mock Services

```go
// Create mock services
mockEmailService := &MockEmailService{}
mockSlackService := &MockSlackService{}

// Create consumer manager with mocks
config := consumers.ConsumerConfig{
    EmailWorkerCount: 1,
    SlackWorkerCount: 1,
}
consumerManager := consumers.NewConsumerManagerWithServices(
    mockEmailService,
    mockSlackService,
    mockKafkaService,
    config,
)

// Test processing
// ... test implementation
```

### Integration Testing

```go
// Use real services for integration tests
emailService := email.NewEmailService()
slackService := slack.NewSlackService()
kafkaService, _ := kafka.NewKafkaService()

consumerManager := consumers.NewConsumerManagerWithServices(
    emailService,
    slackService,
    kafkaService,
    config,
)

// Test end-to-end notification processing
// ... test implementation
```

## Conclusion

This refactoring significantly improves the architecture of the consumer manager by:

1. **Enabling dependency injection** for better testability and flexibility
2. **Promoting service reuse** across the application
3. **Maintaining backward compatibility** with existing code
4. **Providing clear migration paths** for enhanced functionality

The changes are designed to be non-breaking while providing significant benefits for new implementations and testing scenarios. 