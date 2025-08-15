// Package services provides all notification service implementations
package services

import (
	"context"

	"github.com/gaurav2721/notification-service/external_services/email"
	"github.com/gaurav2721/notification-service/external_services/inapp"
	"github.com/gaurav2721/notification-service/external_services/slack"
	"github.com/gaurav2721/notification-service/external_services/user"
	"github.com/gaurav2721/notification-service/notification_manager"
)

// Re-export all interfaces and types for convenience
type (
	EmailService        = email.EmailService
	SlackService        = slack.SlackService
	InAppService        = inapp.InAppService
	UserService         = user.UserService
	NotificationManager = notification_manager.NotificationManager
)

// Re-export all configurations
type (
	EmailConfig        = email.EmailConfig
	SlackConfig        = slack.SlackConfig
	InAppConfig        = inapp.InAppConfig
	UserConfig         = user.UserConfig
	NotificationConfig = notification_manager.NotificationConfig
)

// Re-export all errors
var (
	// Email service errors
	ErrEmailSendFailed       = email.ErrEmailSendFailed
	ErrInvalidEmail          = email.ErrInvalidEmail
	ErrEmailTemplateNotFound = email.ErrEmailTemplateNotFound

	// Slack service errors
	ErrSlackSendFailed   = slack.ErrSlackSendFailed
	ErrInvalidChannel    = slack.ErrInvalidChannel
	ErrSlackTokenMissing = slack.ErrSlackTokenMissing

	// InApp service errors
	ErrInAppSendFailed     = inapp.ErrInAppSendFailed
	ErrInAppDeviceToken    = inapp.ErrInAppDeviceToken
	ErrInAppDeviceNotFound = inapp.ErrInAppDeviceNotFound

	// User service errors
	ErrUserNotFound       = user.ErrUserNotFound
	ErrUserAlreadyExists  = user.ErrUserAlreadyExists
	ErrInvalidUserID      = user.ErrInvalidUserID
	ErrDeviceInactive     = user.ErrDeviceInactive
	ErrInvalidDeviceToken = user.ErrInvalidDeviceToken

	// Notification service errors
	ErrUnsupportedNotificationType = notification_manager.ErrUnsupportedNotificationType
	ErrNoScheduledTime             = notification_manager.ErrNoScheduledTime
	ErrTemplateNotFound            = notification_manager.ErrTemplateNotFound
	ErrInvalidRecipients           = notification_manager.ErrInvalidRecipients
)

// ServiceFactory provides methods to create service instances
type ServiceFactory struct{}

// NewServiceFactory creates a new service factory
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
}

// NewEmailService creates a new email service instance
func (f *ServiceFactory) NewEmailService() EmailService {
	return email.NewEmailService()
}

// NewSlackService creates a new slack service instance
func (f *ServiceFactory) NewSlackService() SlackService {
	return slack.NewSlackService()
}

// NewInAppService creates a new in-app service instance
func (f *ServiceFactory) NewInAppService() InAppService {
	return inapp.NewInAppService()
}

// NewUserService creates a new user service instance
func (f *ServiceFactory) NewUserService() UserService {
	return user.NewUserService()
}

// NewNotificationManager creates a new notification manager instance
func (f *ServiceFactory) NewNotificationManager(
	emailService EmailService,
	slackService SlackService,
	inAppService InAppService,
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(emailService, slackService, inAppService, nil, nil)
}

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
