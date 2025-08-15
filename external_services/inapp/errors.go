package inapp

import "errors"

// InApp service errors
var (
	ErrInAppSendFailed     = errors.New("failed to send in-app notification")
	ErrInAppDeviceToken    = errors.New("invalid in-app device token")
	ErrInAppDeviceNotFound = errors.New("in-app device not found")
)
