// Package services provides all notification service implementations
package services

import (
	"github.com/gaurav2721/notification-service/services/email"
	"github.com/gaurav2721/notification-service/services/inapp"
	"github.com/gaurav2721/notification-service/services/notification"
	"github.com/gaurav2721/notification-service/services/scheduler"
	"github.com/gaurav2721/notification-service/services/slack"
	"github.com/gaurav2721/notification-service/services/user"
)

// Re-export all interfaces and types for convenience
type (
	EmailService        = email.EmailService
	SlackService        = slack.SlackService
	InAppService        = inapp.InAppService
	SchedulerService    = scheduler.SchedulerService
	UserService         = user.UserService
	NotificationService = notification.NotificationService
)

// Re-export all configurations
type (
	EmailConfig        = email.EmailConfig
	SlackConfig        = slack.SlackConfig
	InAppConfig        = inapp.InAppConfig
	SchedulerConfig    = scheduler.SchedulerConfig
	UserConfig         = user.UserConfig
	NotificationConfig = notification.NotificationConfig
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

	// Scheduler service errors
	ErrSchedulingFailed = scheduler.ErrSchedulingFailed
	ErrJobNotFound      = scheduler.ErrJobNotFound
	ErrJobTimeout       = scheduler.ErrJobTimeout

	// User service errors
	ErrUserNotFound       = user.ErrUserNotFound
	ErrUserAlreadyExists  = user.ErrUserAlreadyExists
	ErrInvalidUserID      = user.ErrInvalidUserID
	ErrDeviceInactive     = user.ErrDeviceInactive
	ErrInvalidDeviceToken = user.ErrInvalidDeviceToken

	// Notification service errors
	ErrUnsupportedNotificationType = notification.ErrUnsupportedNotificationType
	ErrNoScheduledTime             = notification.ErrNoScheduledTime
	ErrTemplateNotFound            = notification.ErrTemplateNotFound
	ErrInvalidRecipients           = notification.ErrInvalidRecipients
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

// NewSchedulerService creates a new scheduler service instance
func (f *ServiceFactory) NewSchedulerService() SchedulerService {
	return scheduler.NewSchedulerService()
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
	scheduler SchedulerService,
) NotificationService {
	return notification.NewNotificationManager(emailService, slackService, inAppService, scheduler)
}
