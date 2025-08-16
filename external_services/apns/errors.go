package apns

import "errors"

var (
	// ErrAPNSSendFailed indicates that sending push notification failed
	ErrAPNSSendFailed = errors.New("failed to send APNS push notification")

	// ErrInvalidDeviceToken indicates that the device token is invalid
	ErrInvalidDeviceToken = errors.New("invalid device token")

	// ErrInvalidConfiguration indicates that APNS configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid APNS configuration")

	// ErrDeviceTokenNotFound indicates that device token was not found
	ErrDeviceTokenNotFound = errors.New("device token not found")

	// ErrUserNotFound indicates that user was not found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidNotificationPayload indicates that notification payload is invalid
	ErrInvalidNotificationPayload = errors.New("invalid notification payload")
)
