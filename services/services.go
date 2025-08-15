// Package services provides all notification service implementations
package services

import (
	"github.com/gaurav2721/notification-service/services/common"
	"github.com/gaurav2721/notification-service/services/email"
	"github.com/gaurav2721/notification-service/services/inapp"
	"github.com/gaurav2721/notification-service/services/notification"
	"github.com/gaurav2721/notification-service/services/scheduler"
	"github.com/gaurav2721/notification-service/services/slack"
	"github.com/gaurav2721/notification-service/services/user"
)

// Re-export all interfaces and types for convenience
type (
	EmailService        = common.EmailService
	SlackService        = common.SlackService
	InAppService        = common.InAppService
	SchedulerService    = common.SchedulerService
	UserService         = user.UserService
	NotificationService = notification.NotificationService
)

// Re-export all errors
var (
	ErrUnsupportedNotificationType = common.ErrUnsupportedNotificationType
	ErrNoScheduledTime             = common.ErrNoScheduledTime
	ErrTemplateNotFound            = common.ErrTemplateNotFound
	ErrInvalidRecipients           = common.ErrInvalidRecipients
	ErrEmailSendFailed             = common.ErrEmailSendFailed
	ErrSlackSendFailed             = common.ErrSlackSendFailed
	ErrInAppSendFailed             = common.ErrInAppSendFailed
	ErrSchedulingFailed            = common.ErrSchedulingFailed
)

// Service constructors
func NewEmailService() *email.EmailServiceImpl {
	return email.NewEmailService()
}

func NewSlackService() *slack.SlackServiceImpl {
	return slack.NewSlackService()
}

func NewInAppService() *inapp.InAppServiceImpl {
	return inapp.NewInAppService()
}

func NewSchedulerService() *scheduler.SchedulerServiceImpl {
	return scheduler.NewSchedulerService()
}

func NewUserService() user.UserService {
	return user.NewUserService()
}

func NewNotificationManager(emailService common.EmailService, slackService common.SlackService, inAppService common.InAppService, scheduler common.SchedulerService) *notification.NotificationManager {
	return notification.NewNotificationManager(emailService, slackService, inAppService, scheduler)
}
