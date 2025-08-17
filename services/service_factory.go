// Package services provides all notification service implementations
package services

import (
	"github.com/gaurav2721/notification-service/external_services/apns"
	"github.com/gaurav2721/notification-service/external_services/consumers"
	"github.com/gaurav2721/notification-service/external_services/email"
	"github.com/gaurav2721/notification-service/external_services/fcm"
	"github.com/gaurav2721/notification-service/external_services/kafka"
	"github.com/gaurav2721/notification-service/external_services/slack"
	"github.com/gaurav2721/notification-service/external_services/user"
	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager"
)

// Re-export all interfaces and types for convenience
type (
	EmailService        = email.EmailService
	SlackService        = slack.SlackService
	APNSService         = apns.APNSService
	FCMService          = fcm.FCMService
	UserService         = user.UserService
	KafkaService        = kafka.KafkaService
	ConsumerManager     = consumers.ConsumerManager
	NotificationManager = notification_manager.NotificationManager
)

// Re-export all configurations
type (
	EmailConfig    = email.EmailConfig
	SlackConfig    = slack.SlackConfig
	APNSConfig     = apns.APNSConfig
	FCMConfig      = fcm.FCMConfig
	ConsumerConfig = consumers.ConsumerConfig
)

// Re-export all errors
var (
	// Email service errors
	ErrEmailSendFailed       = email.ErrEmailSendFailed
	ErrInvalidEmail          = models.ErrInvalidEmail
	ErrEmailTemplateNotFound = email.ErrEmailTemplateNotFound

	// Slack service errors
	ErrSlackSendFailed   = slack.ErrSlackSendFailed
	ErrInvalidChannel    = slack.ErrInvalidChannel
	ErrSlackTokenMissing = slack.ErrSlackTokenMissing

	// APNS service errors
	ErrAPNSSendFailed                 = apns.ErrAPNSSendFailed
	ErrAPNSInvalidConfiguration       = apns.ErrInvalidConfiguration
	ErrAPNSInvalidNotificationPayload = apns.ErrInvalidNotificationPayload

	// FCM service errors
	ErrFCMSendFailed       = fcm.ErrFCMSendFailed
	ErrFCMInvalidConfig    = fcm.ErrInvalidConfiguration
	ErrFCMInvalidPayload   = fcm.ErrInvalidNotificationPayload
	ErrFCMInvalidServerKey = fcm.ErrInvalidServerKey

	// User service errors
	ErrUserNotFound      = user.ErrUserNotFound
	ErrUserAlreadyExists = user.ErrUserAlreadyExists
	ErrInvalidUserID     = user.ErrInvalidUserID
	ErrDeviceInactive    = user.ErrDeviceInactive

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

// NewAPNSService creates a new APNS service instance
func (f *ServiceFactory) NewAPNSService() APNSService {
	return apns.NewAPNSService()
}

// NewFCMService creates a new FCM service instance
func (f *ServiceFactory) NewFCMService() FCMService {
	return fcm.NewFCMService()
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

// NewNotificationManager creates a new notification manager instance
// The scheduler is initialized internally within the notification manager
func (f *ServiceFactory) NewNotificationManager(
	userService UserService,
	kafkaService KafkaService,
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(userService, kafkaService)
}

// NewNotificationManagerWithScheduler creates a new notification manager
// The scheduler is initialized internally within the notification manager
func (f *ServiceFactory) NewNotificationManagerWithScheduler(
	userService UserService,
	kafkaService KafkaService,
) NotificationManager {
	return notification_manager.NewNotificationManagerWithDefaultTemplate(userService, kafkaService)
}

// Note: Scheduler is now initialized internally within the notification manager
// No need to create it externally
