package fcm

import "errors"

var (
	// ErrFCMSendFailed indicates that sending push notification failed
	ErrFCMSendFailed = errors.New("failed to send FCM push notification")

	// ErrInvalidConfiguration indicates that FCM configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid FCM configuration")

	// ErrInvalidNotificationPayload indicates that notification payload is invalid
	ErrInvalidNotificationPayload = errors.New("invalid notification payload")

	// ErrInvalidServerKey indicates that FCM server key is invalid
	ErrInvalidServerKey = errors.New("invalid FCM server key")
)
