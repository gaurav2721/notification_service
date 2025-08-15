// Package services provides all notification service implementations
package services

import (
	"context"

	"github.com/gaurav2721/notification-service/external_services/consumers"
	"github.com/gaurav2721/notification-service/external_services/kafka"
)

// ServiceContainer manages all service dependencies
type ServiceContainer struct {
	emailService        EmailService
	slackService        SlackService
	inAppService        InAppService
	userService         UserService
	kafkaService        kafka.KafkaService
	consumerManager     consumers.ConsumerManager
	notificationService NotificationManager
}

// NewServiceContainer creates a new service container with all dependencies
func NewServiceContainer() *ServiceContainer {
	container := &ServiceContainer{}
	container.initializeServices()
	return container
}

// initializeServices sets up all service dependencies
func (c *ServiceContainer) initializeServices() {
	// Create service factory
	factory := NewServiceFactory()

	// Initialize core services
	c.emailService = factory.NewEmailService()
	c.slackService = factory.NewSlackService()
	c.inAppService = factory.NewInAppService()
	c.userService = factory.NewUserService()

	// Initialize Kafka service using factory
	kafkaService, err := factory.NewKafkaService()
	if err != nil {
		panic("Failed to initialize Kafka service: " + err.Error())
	}
	c.kafkaService = kafkaService

	// Initialize consumer manager using factory with environment configuration
	c.consumerManager = factory.NewConsumerManagerFromEnv(c.kafkaService)

	// Start the consumer manager immediately
	ctx := context.Background()
	if err := c.consumerManager.Initialize(ctx); err != nil {
		panic("Failed to initialize consumer manager: " + err.Error())
	}
	if err := c.consumerManager.Start(ctx); err != nil {
		panic("Failed to start consumer manager: " + err.Error())
	}

	// Initialize notification service with user service and Kafka service
	c.notificationService = factory.NewNotificationManagerWithUserService(c.userService, c.kafkaService)
}

// GetEmailService returns the email service
func (c *ServiceContainer) GetEmailService() EmailService {
	return c.emailService
}

// GetSlackService returns the slack service
func (c *ServiceContainer) GetSlackService() SlackService {
	return c.slackService
}

// GetInAppService returns the in-app service
func (c *ServiceContainer) GetInAppService() InAppService {
	return c.inAppService
}

// GetUserService returns the user service
func (c *ServiceContainer) GetUserService() UserService {
	return c.userService
}

// GetKafkaService returns the kafka service
func (c *ServiceContainer) GetKafkaService() kafka.KafkaService {
	return c.kafkaService
}

// GetConsumerManager returns the consumer manager
func (c *ServiceContainer) GetConsumerManager() consumers.ConsumerManager {
	return c.consumerManager
}

// GetNotificationService returns the notification service
func (c *ServiceContainer) GetNotificationService() NotificationManager {
	return c.notificationService
}

// Shutdown gracefully shuts down all services
func (c *ServiceContainer) Shutdown(ctx context.Context) error {
	// Stop consumer manager
	if c.consumerManager != nil {
		if err := c.consumerManager.Stop(); err != nil {
			// Log error but continue with shutdown
			// TODO: Add proper logging
		}
	}

	// Close Kafka service
	if c.kafkaService != nil {
		c.kafkaService.Close()
	}

	// Add any cleanup logic here if needed
	// For now, all services are stateless, so no cleanup is required
	return nil
}

// ServiceProvider interface for dependency injection
type ServiceProvider interface {
	GetEmailService() EmailService
	GetSlackService() SlackService
	GetInAppService() InAppService
	GetUserService() UserService
	GetKafkaService() kafka.KafkaService
	GetConsumerManager() consumers.ConsumerManager
	GetNotificationService() NotificationManager
	Shutdown(ctx context.Context) error
}

// Ensure ServiceContainer implements ServiceProvider
var _ ServiceProvider = (*ServiceContainer)(nil)
