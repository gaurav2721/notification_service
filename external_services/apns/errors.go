package apns

import "errors"

var (
	// ErrAPNSSendFailed indicates that sending push notification failed
	ErrAPNSSendFailed = errors.New("failed to send APNS push notification")

	// ErrInvalidConfiguration indicates that APNS configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid APNS configuration")

	// ErrInvalidNotificationPayload indicates that notification payload is invalid
	ErrInvalidNotificationPayload = errors.New("invalid notification payload")
)
