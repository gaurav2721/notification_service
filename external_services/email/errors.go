package email

import "errors"

// Email service errors
var (
	ErrEmailSendFailed       = errors.New("failed to send email")
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrEmailTemplateNotFound = errors.New("email template not found")
)
