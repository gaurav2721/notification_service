package inapp

import "context"

// InAppService interface defines methods for in-app notifications
type InAppService interface {
	SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error)
}
