# Consumer Manager Refactoring: Complete Service Dependencies

This document explains the complete refactoring of the consumer manager to accept all service dependencies including email, slack, APNS, and FCM services.

## Overview

The consumer manager has been fully refactored to support dependency injection of all notification services, enabling comprehensive testability, flexibility, and service reuse across the application.

## Key Changes

### 1. Updated ConsumerConfig Interface

**File**: `external_services/consumers/interfaces.go`

Added all service dependencies to the `ConsumerConfig` struct:

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
    APNSService interface {
        SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
    }
    FCMService interface {
        SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
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

### 2. Enhanced Consumer Manager Constructor

**File**: `external_services/consumers/manager.go`

Updated constructor to accept all service dependencies:

```go
// NewConsumerManagerWithServices creates a new consumer manager with service dependencies
func NewConsumerManagerWithServices(
    emailService interface {
        SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
    },
    slackService interface {
        SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
    },
    apnsService interface {
        SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
    },
    fcmService interface {
        SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
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
    config.APNSService = apnsService
    config.FCMService = fcmService
    config.KafkaService = kafkaService

    return &consumerManager{
        config:      config,
        workerPools: make(map[NotificationType]ConsumerWorkerPool),
        running:     false,
    }
}
```

### 3. Updated All Processor Creation Methods

**File**: `external_services/consumers/manager.go`

All processor creation methods now use injected services:

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

// createIOSPushWorkerPool creates the iOS push notification worker pool
func (cm *consumerManager) createIOSPushWorkerPool() {
    var processor NotificationProcessor
    
    // Use injected APNS service if available, otherwise create default
    if cm.config.APNSService != nil {
        processor = NewIOSPushProcessorWithService(cm.config.APNSService)
    } else {
        processor = NewIOSPushProcessor()
    }
    
    pool := NewWorkerPool(
        IOSPushNotification,
        cm.config.KafkaService.GetIOSPushNotificationChannel(),
        processor,
        cm.config.IOSPushWorkerCount,
    )
    cm.workerPools[IOSPushNotification] = pool
}

// createAndroidPushWorkerPool creates the Android push notification worker pool
func (cm *consumerManager) createAndroidPushWorkerPool() {
    var processor NotificationProcessor
    
    // Use injected FCM service if available, otherwise create default
    if cm.config.FCMService != nil {
        processor = NewAndroidPushProcessorWithService(cm.config.FCMService)
    } else {
        processor = NewAndroidPushProcessor()
    }
    
    pool := NewWorkerPool(
        AndroidPushNotification,
        cm.config.KafkaService.GetAndroidPushNotificationChannel(),
        processor,
        cm.config.AndroidPushWorkerCount,
    )
    cm.workerPools[AndroidPushNotification] = pool
}
```

### 4. Enhanced All Processors

#### Email Processor
**File**: `external_services/consumers/email_processor.go`

Already supports service injection with `NewEmailProcessorWithService()`.

#### Slack Processor
**File**: `external_services/consumers/slack_processor.go`

Enhanced with service injection and actual processing logic:

```go
type slackProcessor struct {
    slackService interface {
        SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
    }
}

func NewSlackProcessorWithService(slackService interface {
    SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
    return &slackProcessor{
        slackService: slackService,
    }
}
```

#### iOS Push Processor
**File**: `external_services/consumers/ios_push_processor.go`

Enhanced with APNS service injection and actual processing logic:

```go
type iosPushProcessor struct {
    apnsService interface {
        SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
    }
}

func NewIOSPushProcessorWithService(apnsService interface {
    SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
    return &iosPushProcessor{
        apnsService: apnsService,
    }
}
```

#### Android Push Processor
**File**: `external_services/consumers/android_push_processor.go`

Enhanced with FCM service injection and actual processing logic:

```go
type androidPushProcessor struct {
    fcmService interface {
        SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
    }
}

func NewAndroidPushProcessorWithService(fcmService interface {
    SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
}) NotificationProcessor {
    return &androidPushProcessor{
        fcmService: fcmService,
    }
}
```

### 5. Updated Service Provider

**File**: `services/service_provider.go`

Modified to use the new constructor with all service dependencies:

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
    c.apnsService,
    c.fcmService,
    c.kafkaService,
    config,
)
```

### 6. Enhanced Service Factory

**File**: `services/service_factory.go`

Updated factory method to accept all services:

```go
// NewConsumerManagerWithServices creates a new consumer manager with service dependencies
func (f *ServiceFactory) NewConsumerManagerWithServices(
    emailService EmailService,
    slackService SlackService,
    apnsService APNSService,
    fcmService FCMService,
    kafkaService KafkaService,
    config ConsumerConfig,
) ConsumerManager {
    return consumers.NewConsumerManagerWithServices(emailService, slackService, apnsService, fcmService, kafkaService, config)
}
```

## Benefits of Complete Refactoring

### 1. **Comprehensive Dependency Injection**
- All services can now be injected into the consumer manager
- Complete testability with mock services
- Flexible service configuration for all notification types

### 2. **Service Reuse Across Application**
- The same services can be used across the entire application
- Consistent service behavior and configuration
- Reduced resource usage and improved performance

### 3. **Enhanced Testing Capabilities**
- Easy to mock all services for unit tests
- Isolated testing of all processors
- Better integration test scenarios
- Comprehensive test coverage

### 4. **Maximum Flexibility**
- Can use different service implementations for each type
- Easy to switch between real and mock services
- Support for custom service implementations
- Platform-specific optimizations

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
    IOSPushWorkerCount: 3,
    AndroidPushWorkerCount: 3,
}
consumerManager := consumers.NewConsumerManager(config)
```

### Advanced Usage with All Custom Services

```go
// Create custom services
emailService := &CustomEmailService{}
slackService := &CustomSlackService{}
apnsService := &CustomAPNSService{}
fcmService := &CustomFCMService{}
kafkaService := &KafkaService{}

// Create consumer manager with all custom services
config := consumers.ConsumerConfig{
    EmailWorkerCount: 5,
    SlackWorkerCount: 3,
    IOSPushWorkerCount: 3,
    AndroidPushWorkerCount: 3,
}
consumerManager := consumers.NewConsumerManagerWithServices(
    emailService,
    slackService,
    apnsService,
    fcmService,
    kafkaService,
    config,
)
```

### Using Service Factory

```go
factory := services.NewServiceFactory()

// Create all services
emailService := factory.NewEmailService()
slackService := factory.NewSlackService()
apnsService := factory.NewAPNSService()
fcmService := factory.NewFCMService()
kafkaService, _ := factory.NewKafkaService()

// Create consumer manager with all services
config := consumers.ConsumerConfig{
    EmailWorkerCount: 5,
    SlackWorkerCount: 3,
    IOSPushWorkerCount: 3,
    AndroidPushWorkerCount: 3,
}
consumerManager := factory.NewConsumerManagerWithServices(
    emailService,
    slackService,
    apnsService,
    fcmService,
    kafkaService,
    config,
)
```

## Notification Payload Examples

### Email Notification
```json
{
  "notification_id": "123e4567-e89b-12d3-a456-426614174000",
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
}
```

### Slack Notification
```json
{
  "notification_id": "456e7890-e89b-12d3-a456-426614174001",
  "type": "slack",
  "content": {
    "text": "Hello from our service!"
  },
  "channel": "#general"
}
```

### iOS Push Notification
```json
{
  "notification_id": "789e0123-e89b-12d3-a456-426614174002",
  "type": "ios_push",
  "content": {
    "title": "New Message",
    "body": "You have a new message"
  },
  "recipients": ["ios_device_token_123", "ios_device_token_456"]
}
```

### Android Push Notification
```json
{
  "notification_id": "101e2345-e89b-12d3-a456-426614174003",
  "type": "android_push",
  "content": {
    "title": "New Message",
    "body": "You have a new message"
  },
  "recipients": ["android_device_token_123", "android_device_token_456"]
}
```

## Testing Examples

### Unit Testing with Mock Services

```go
// Create mock services
mockEmailService := &MockEmailService{}
mockSlackService := &MockSlackService{}
mockAPNSService := &MockAPNSService{}
mockFCMService := &MockFCMService{}

// Create consumer manager with mocks
config := consumers.ConsumerConfig{
    EmailWorkerCount: 1,
    SlackWorkerCount: 1,
    IOSPushWorkerCount: 1,
    AndroidPushWorkerCount: 1,
}
consumerManager := consumers.NewConsumerManagerWithServices(
    mockEmailService,
    mockSlackService,
    mockAPNSService,
    mockFCMService,
    mockKafkaService,
    config,
)

// Test processing for all notification types
// ... test implementation
```

### Integration Testing

```go
// Use real services for integration tests
emailService := email.NewEmailService()
slackService := slack.NewSlackService()
apnsService := apns.NewAPNSService()
fcmService := fcm.NewFCMService()
kafkaService, _ := kafka.NewKafkaService()

consumerManager := consumers.NewConsumerManagerWithServices(
    emailService,
    slackService,
    apnsService,
    fcmService,
    kafkaService,
    config,
)

// Test end-to-end notification processing for all types
// ... test implementation
```

## Migration Guide

### For Existing Code

1. **No Changes Required**: Existing code using `NewConsumerManager()` continues to work
2. **Optional Enhancement**: Update to use `NewConsumerManagerWithServices()` for better dependency management
3. **Service Provider**: The service provider automatically uses the new constructor with all services

### For New Code

1. **Use New Constructor**: Prefer `NewConsumerManagerWithServices()` for new implementations
2. **Inject All Services**: Pass all required services as dependencies
3. **Use Factory**: Leverage the service factory for consistent service creation
4. **Test with Mocks**: Use mock services for comprehensive testing

## Conclusion

This complete refactoring significantly improves the architecture of the consumer manager by:

1. **Enabling comprehensive dependency injection** for all notification services
2. **Promoting service reuse** across the entire application
3. **Maintaining backward compatibility** with existing code
4. **Providing clear migration paths** for enhanced functionality
5. **Supporting all notification types** with consistent patterns

The changes are designed to be non-breaking while providing significant benefits for new implementations, testing scenarios, and service management across the application. 