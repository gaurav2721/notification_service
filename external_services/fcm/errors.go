package fcm

import "errors"

var (
	// ErrFCMSendFailed indicates that sending push notification failed
	ErrFCMSendFailed = errors.New("failed to send FCM push notification")

	// ErrInvalidDeviceToken indicates that the device token is invalid
	ErrInvalidDeviceToken = errors.New("invalid device token")

	// ErrInvalidConfiguration indicates that FCM configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid FCM configuration")

	// ErrDeviceTokenNotFound indicates that device token was not found
	ErrDeviceTokenNotFound = errors.New("device token not found")

	// ErrUserNotFound indicates that user was not found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidNotificationPayload indicates that notification payload is invalid
	ErrInvalidNotificationPayload = errors.New("invalid notification payload")

	// ErrInvalidServerKey indicates that FCM server key is invalid
	ErrInvalidServerKey = errors.New("invalid FCM server key")
)
