package fcm

import "context"

// FCMService interface defines methods for Firebase Cloud Messaging
type FCMService interface {
	SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
}

// FCMConfig holds configuration for FCM service
type FCMConfig struct {
	ServerKey string
	ProjectID string
	Timeout   int // in seconds
	BatchSize int // number of tokens to send in a single request
}
