// Package services provides all notification service implementations
package services

import (
	"context"

	"github.com/gaurav2721/notification-service/external_services/consumers"
	"github.com/gaurav2721/notification-service/external_services/kafka"
	"github.com/sirupsen/logrus"
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
	logrus.Debug("Creating new service container")
	container := &ServiceContainer{}
	container.initializeServices()
	logrus.Debug("Service container created successfully")
	return container
}

// initializeServices sets up all service dependencies
func (c *ServiceContainer) initializeServices() {
	logrus.Debug("Initializing service dependencies")

	// Create service factory
	factory := NewServiceFactory()
	logrus.Debug("Service factory created")

	// Initialize core services
	logrus.Debug("Initializing core services")
	c.emailService = factory.NewEmailService()
	c.slackService = factory.NewSlackService()
	c.inAppService = factory.NewInAppService()
	c.userService = factory.NewUserService()
	logrus.Debug("Core services initialized")

	// Initialize Kafka service using factory
	logrus.Debug("Initializing Kafka service")
	kafkaService, err := factory.NewKafkaService()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize Kafka service")
		panic("Failed to initialize Kafka service: " + err.Error())
	}
	c.kafkaService = kafkaService
	logrus.Debug("Kafka service initialized successfully")

	// Initialize consumer manager using factory with environment configuration
	logrus.Debug("Initializing consumer manager")
	c.consumerManager = factory.NewConsumerManagerFromEnv(c.kafkaService)

	// Start the consumer manager immediately
	ctx := context.Background()
	logrus.Debug("Starting consumer manager")
	if err := c.consumerManager.Initialize(ctx); err != nil {
		logrus.WithError(err).Fatal("Failed to initialize consumer manager")
		panic("Failed to initialize consumer manager: " + err.Error())
	}
	if err := c.consumerManager.Start(ctx); err != nil {
		logrus.WithError(err).Fatal("Failed to start consumer manager")
		panic("Failed to start consumer manager: " + err.Error())
	}
	logrus.Debug("Consumer manager started successfully")

	// Initialize notification service with user service and Kafka service
	logrus.Debug("Initializing notification service")
	c.notificationService = factory.NewNotificationManagerWithUserService(c.userService, c.kafkaService)
	logrus.Debug("Notification service initialized")

	logrus.Debug("All service dependencies initialized successfully")
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
	logrus.Debug("Starting graceful shutdown of service container")

	// Stop consumer manager
	if c.consumerManager != nil {
		logrus.Debug("Stopping consumer manager")
		if err := c.consumerManager.Stop(); err != nil {
			logrus.WithError(err).Error("Error stopping consumer manager")
		} else {
			logrus.Debug("Consumer manager stopped successfully")
		}
	}

	// Close Kafka service
	if c.kafkaService != nil {
		logrus.Debug("Closing Kafka service")
		c.kafkaService.Close()
		logrus.Debug("Kafka service closed")
	}

	// Add any cleanup logic here if needed
	// For now, all services are stateless, so no cleanup is required
	logrus.Debug("Service container shutdown completed")
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
