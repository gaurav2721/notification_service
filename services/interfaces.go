package services

import (
	"context"
	"errors"
	"time"
)

// Service interfaces
type EmailService interface {
	SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
}

type SlackService interface {
	SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
}

type InAppService interface {
	SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error)
}

type SchedulerService interface {
	ScheduleJob(jobID string, scheduledTime time.Time, job func()) error
	CancelJob(jobID string) error
}

// Custom errors
var (
	ErrUnsupportedNotificationType = errors.New("unsupported notification type")
	ErrNoScheduledTime             = errors.New("no scheduled time provided")
	ErrTemplateNotFound            = errors.New("template not found")
	ErrInvalidRecipients           = errors.New("invalid recipients")
	ErrEmailSendFailed             = errors.New("failed to send email")
	ErrSlackSendFailed             = errors.New("failed to send slack message")
	ErrInAppSendFailed             = errors.New("failed to send in-app notification")
	ErrSchedulingFailed            = errors.New("failed to schedule notification")
)
