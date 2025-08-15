# Service Architecture Documentation

## Overview

The notification service follows a clean architecture pattern with strict separation of concerns and dependency inversion principles. All services are hidden behind interfaces, making the system highly testable, maintainable, and extensible.

## Architecture Principles

### 1. Interface Segregation Principle
- Each service has a well-defined interface
- Interfaces are small and focused on specific responsibilities
- No service depends on concrete implementations

### 2. Dependency Inversion Principle
- High-level modules (handlers) depend on abstractions (interfaces)
- Low-level modules (implementations) implement abstractions
- Dependencies are injected, not created internally

### 3. Single Responsibility Principle
- Each service has a single, well-defined responsibility
- Services are loosely coupled
- Changes to one service don't affect others

## Service Layer Structure

```
services/
├── services.go          # Main service interfaces and factory
├── container.go         # Service container for dependency management
├── config.go           # Configuration structures and factories
├── common/             # Shared interfaces and types
├── user/               # User service implementation
├── email/              # Email service implementation
├── slack/              # Slack service implementation
├── inapp/              # In-app notification service implementation
├── notification/       # Notification manager implementation
└── scheduler/          # Scheduler service implementation
```

## Core Interfaces

### ServiceProvider Interface
```go
type ServiceProvider interface {
    GetEmailService() EmailService
    GetSlackService() SlackService
    GetInAppService() InAppService
    GetSchedulerService() SchedulerService
    GetUserService() UserService
    GetNotificationService() NotificationService
    Shutdown(ctx context.Context) error
}
```

### Individual Service Interfaces
- `EmailService`: Handles email notifications
- `SlackService`: Handles Slack notifications
- `InAppService`: Handles in-app notifications
- `SchedulerService`: Handles notification scheduling
- `UserService`: Handles user and device management
- `NotificationService`: Orchestrates all notification types

## Service Factory Pattern

### Basic Factory
```go
// Create a basic service factory
factory := services.NewServiceFactory()

// Create services (returns interfaces)
emailService := factory.NewEmailService()
userService := factory.NewUserService()
```

### Configurable Factory
```go
// Create configuration
config := services.DefaultServiceConfig()
config.EmailConfig.SMTPHost = "smtp.company.com"

// Create factory with configuration
factory := services.NewServiceFactoryWithConfig(config)
emailService := factory.NewEmailService()
```

## Service Container Pattern

### Basic Container
```go
// Create service container (manages all dependencies)
container := services.NewServiceContainer()

// Get services from container
userService := container.GetUserService()
notificationService := container.GetNotificationService()
```

### Configurable Container
```go
// Create configuration
config := services.DefaultServiceConfig()
config.UserConfig.DatabaseURL = "postgres://localhost/users"

// Create container with configuration
container := services.NewServiceContainerWithConfig(config)
userService := container.GetUserService()
```

## Dependency Injection in Main

### Current Implementation
```go
func main() {
    // Initialize service container
    serviceContainer := services.NewServiceContainer()

    // Initialize handlers with interface dependencies
    notificationHandler := handlers.NewNotificationHandler(
        serviceContainer.GetNotificationService(),
    )
    userHandler := handlers.NewUserHandler(
        serviceContainer.GetUserService(),
    )

    // Setup routes and start server
    // ...
}
```

## Benefits of This Architecture

### 1. Testability
- Easy to mock services for unit testing
- Handlers can be tested in isolation
- Integration tests can use real or mock services

### 2. Maintainability
- Changes to service implementations don't affect consumers
- New implementations can be added without changing existing code
- Clear separation of concerns

### 3. Extensibility
- New notification channels can be added easily
- Different storage backends can be swapped
- Configuration can be changed without code changes

### 4. Flexibility
- Services can be configured differently for different environments
- Multiple implementations can coexist
- Easy to add middleware or decorators

## Testing Strategy

### Unit Testing
```go
// Mock service for testing
type MockUserService struct {
    users map[string]*models.User
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
    if user, exists := m.users[id]; exists {
        return user, nil
    }
    return nil, user.ErrUserNotFound
}

// Test handler with mock
func TestUserHandler_GetUser(t *testing.T) {
    mockService := &MockUserService{
        users: map[string]*models.User{
            "user-001": &models.User{ID: "user-001", Email: "test@example.com"},
        },
    }
    
    handler := handlers.NewUserHandler(mockService)
    // Test implementation...
}
```

### Integration Testing
```go
func TestUserServiceIntegration(t *testing.T) {
    // Use real service with test configuration
    config := services.DefaultServiceConfig()
    config.UserConfig.DatabaseURL = "memory://test"
    
    container := services.NewServiceContainerWithConfig(config)
    userService := container.GetUserService()
    
    // Test with real implementation...
}
```

## Configuration Management

### Environment-Based Configuration
```go
func loadConfig() *services.ServiceConfig {
    config := services.DefaultServiceConfig()
    
    // Override with environment variables
    if smtpHost := os.Getenv("SMTP_HOST"); smtpHost != "" {
        config.EmailConfig.SMTPHost = smtpHost
    }
    
    if slackToken := os.Getenv("SLACK_BOT_TOKEN"); slackToken != "" {
        config.SlackConfig.BotToken = slackToken
    }
    
    return config
}
```

### Configuration Validation
```go
func (c *ServiceConfig) Validate() error {
    if c.EmailConfig.SMTPHost == "" {
        return errors.New("SMTP host is required")
    }
    
    if c.SlackConfig.BotToken == "" {
        return errors.New("Slack bot token is required")
    }
    
    return nil
}
```

## Migration Guide

### From Direct Service Usage
**Before:**
```go
// Direct instantiation
userService := user.NewUserService()
handler := handlers.NewUserHandler(userService)
```

**After:**
```go
// Interface-based dependency injection
container := services.NewServiceContainer()
handler := handlers.NewUserHandler(container.GetUserService())
```

### From Concrete Types
**Before:**
```go
func NewHandler(userService *user.UserServiceImpl) *Handler {
    return &Handler{userService: userService}
}
```

**After:**
```go
func NewHandler(userService user.UserService) *Handler {
    return &Handler{userService: userService}
}
```

## Best Practices

### 1. Always Use Interfaces
- Never depend on concrete implementations
- Define interfaces at the package level
- Keep interfaces small and focused

### 2. Use Dependency Injection
- Inject dependencies through constructors
- Use service containers for complex dependency graphs
- Avoid global state and singletons

### 3. Configuration Management
- Use configuration structures for service setup
- Validate configuration at startup
- Support environment-based configuration

### 4. Error Handling
- Define service-specific errors
- Use error wrapping for context
- Provide meaningful error messages

### 5. Testing
- Write unit tests for all services
- Use mocks for external dependencies
- Test error conditions and edge cases

## Future Enhancements

### 1. Service Discovery
- Add service registry for dynamic service discovery
- Support for service health checks
- Load balancing between service instances

### 2. Observability
- Add metrics collection
- Structured logging
- Distributed tracing

### 3. Caching
- Add caching layer for frequently accessed data
- Cache invalidation strategies
- Multi-level caching

### 4. Circuit Breakers
- Add circuit breaker pattern for external services
- Fallback mechanisms
- Graceful degradation

This architecture provides a solid foundation for building scalable, maintainable, and testable notification services while following Go best practices and clean architecture principles. 