package notification

import "errors"

// Notification service errors
var (
	ErrUnsupportedNotificationType = errors.New("unsupported notification type")
	ErrNoScheduledTime             = errors.New("no scheduled time provided")
	ErrTemplateNotFound            = errors.New("template not found")
	ErrInvalidRecipients           = errors.New("invalid recipients")
)
