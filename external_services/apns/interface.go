package apns

import "context"

// APNSService interface defines methods for Apple Push Notification Service
type APNSService interface {
	SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	RegisterDeviceToken(userID, deviceToken string) error
	UnregisterDeviceToken(userID, deviceToken string) error
	GetDeviceTokensForUser(userID string) ([]string, error)
}

// APNSConfig holds configuration for APNS service
type APNSConfig struct {
	BundleID       string
	KeyID          string
	TeamID         string
	PrivateKeyPath string
	Environment    string // "sandbox" or "production"
	MaxRetries     int
	Timeout        int // in seconds
}
