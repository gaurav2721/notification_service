package email

import "context"

// EmailService interface defines methods for email notifications
type EmailService interface {
	SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
}
