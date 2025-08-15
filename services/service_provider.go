// Package services provides all notification service implementations
package services

import (
	"context"
)

// ServiceContainer manages all service dependencies
type ServiceContainer struct {
	emailService        EmailService
	slackService        SlackService
	inAppService        InAppService
	userService         UserService
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

	// Initialize notification service with dependencies
	c.notificationService = factory.NewNotificationManager(
		c.emailService,
		c.slackService,
		c.inAppService,
	)
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

// GetNotificationService returns the notification service
func (c *ServiceContainer) GetNotificationService() NotificationManager {
	return c.notificationService
}

// Shutdown gracefully shuts down all services
func (c *ServiceContainer) Shutdown(ctx context.Context) error {
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
	GetNotificationService() NotificationManager
	Shutdown(ctx context.Context) error
}

// Ensure ServiceContainer implements ServiceProvider
var _ ServiceProvider = (*ServiceContainer)(nil)
