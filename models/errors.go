package models

import "errors"

// Template-related errors
var (
	ErrInvalidTemplateContent  = errors.New("invalid template content")
	ErrInvalidTemplateType     = errors.New("invalid template type")
	ErrMissingRequiredVariable = errors.New("missing required variable")
	ErrTemplateNotFound        = errors.New("template not found")
)

// Email-related errors
var (
	ErrInvalidEmail           = errors.New("invalid email address")
	ErrMissingEmailID         = errors.New("email notification ID is required")
	ErrMissingEmailType       = errors.New("email notification type is required")
	ErrMissingEmailContent    = errors.New("email notification content is required")
	ErrMissingEmailRecipients = errors.New("email notification recipients are required")
	ErrMissingEmailFrom       = errors.New("email notification from email is required")
	ErrMissingEmailSubject    = errors.New("email subject is required")
	ErrMissingEmailBody       = errors.New("email body is required")
	ErrEmptyEmailRecipients   = errors.New("recipients list cannot be empty")
)
