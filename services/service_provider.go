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
	schedulerService    SchedulerService
	userService         UserService
	notificationService NotificationService
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
	c.schedulerService = factory.NewSchedulerService()
	c.userService = factory.NewUserService()

	// Initialize notification service with dependencies
	c.notificationService = factory.NewNotificationManager(
		c.emailService,
		c.slackService,
		c.inAppService,
		c.schedulerService,
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

// GetSchedulerService returns the scheduler service
func (c *ServiceContainer) GetSchedulerService() SchedulerService {
	return c.schedulerService
}

// GetUserService returns the user service
func (c *ServiceContainer) GetUserService() UserService {
	return c.userService
}

// GetNotificationService returns the notification service
func (c *ServiceContainer) GetNotificationService() NotificationService {
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
	GetSchedulerService() SchedulerService
	GetUserService() UserService
	GetNotificationService() NotificationService
	Shutdown(ctx context.Context) error
}

// Ensure ServiceContainer implements ServiceProvider
var _ ServiceProvider = (*ServiceContainer)(nil)
