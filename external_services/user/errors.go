package user

import "errors"

// User service errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrDeviceNotFound     = errors.New("device not found")
	ErrDeviceInactive     = errors.New("device is inactive")
	ErrInvalidDeviceToken = errors.New("invalid device token")
)
