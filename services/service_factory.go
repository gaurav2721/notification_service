// Package services provides all notification service implementations
package services

import (
	"github.com/gaurav2721/notification-service/external_services/consumers"
	"github.com/gaurav2721/notification-service/external_services/email"
	"github.com/gaurav2721/notification-service/external_services/inapp"
	"github.com/gaurav2721/notification-service/external_services/kafka"
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
	KafkaService        = kafka.KafkaService
	ConsumerManager     = consumers.ConsumerManager
	NotificationManager = notification_manager.NotificationManager
)

// Re-export all configurations
type (
	EmailConfig        = email.EmailConfig
	SlackConfig        = slack.SlackConfig
	InAppConfig        = inapp.InAppConfig
	UserConfig         = user.UserConfig
	ConsumerConfig     = consumers.ConsumerConfig
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

// NewKafkaService creates a new kafka service instance
func (f *ServiceFactory) NewKafkaService() (KafkaService, error) {
	return kafka.NewKafkaService()
}

// NewConsumerManager creates a new consumer manager instance
func (f *ServiceFactory) NewConsumerManager(kafkaService KafkaService) ConsumerManager {
	config := consumers.ConsumerConfig{
		KafkaService: kafkaService,
	}
	return consumers.NewConsumerManager(config)
}

// NewConsumerManagerWithConfig creates a new consumer manager with custom configuration
func (f *ServiceFactory) NewConsumerManagerWithConfig(config ConsumerConfig) ConsumerManager {
	return consumers.NewConsumerManager(config)
}

// NewConsumerManagerFromEnv creates a new consumer manager with configuration from environment variables
func (f *ServiceFactory) NewConsumerManagerFromEnv(kafkaService KafkaService) ConsumerManager {
	return consumers.NewConsumerManagerFromEnv(kafkaService)
}

// NewNotificationManager creates a new notification manager instance
func (f *ServiceFactory) NewNotificationManager(
	userService UserService,
	kafkaService KafkaService,
	scheduler interface{},
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(userService, kafkaService, scheduler)
}

// NewNotificationManagerWithUserService creates a new notification manager with user service
func (f *ServiceFactory) NewNotificationManagerWithUserService(
	userService UserService,
	kafkaService KafkaService,
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(userService, kafkaService, nil)
}

// NewNotificationManagerWithScheduler creates a new notification manager with scheduler
func (f *ServiceFactory) NewNotificationManagerWithScheduler(
	userService UserService,
	kafkaService KafkaService,
	scheduler interface{},
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(userService, kafkaService, scheduler)
}

// NewNotificationManagerComplete creates a new notification manager with all dependencies
func (f *ServiceFactory) NewNotificationManagerComplete(
	userService UserService,
	kafkaService KafkaService,
	scheduler interface{},
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(userService, kafkaService, scheduler)
}

// NewNotificationManagerWithKafkaOnly creates a new notification manager with only Kafka service
// This is used when the notification manager only needs to push notifications to Kafka channels
func (f *ServiceFactory) NewNotificationManagerWithKafkaOnly(kafkaService KafkaService) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(nil, kafkaService, nil)
}
