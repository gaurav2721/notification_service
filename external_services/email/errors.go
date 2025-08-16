package email

import "errors"

// Email service errors
var (
	ErrEmailSendFailed       = errors.New("failed to send email")
	ErrEmailTemplateNotFound = errors.New("email template not found")
)
