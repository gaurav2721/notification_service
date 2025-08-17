package user

import "github.com/gaurav2721/notification-service/models"

// User service interface and related types can be added here
// UserService interface defines methods for user management
type UserService interface {
	GetUserByID(userID string) (*models.User, error)
	GetUsersByIDs(userIDs []string) ([]*models.User, error)
	GetAllUsers() ([]*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error

	// Device management methods
	RegisterDevice(userID, deviceToken, deviceType string) (*models.UserDeviceInfo, error)
	GetUserDevices(userID string) ([]*models.UserDeviceInfo, error)
	GetActiveUserDevices(userID string) ([]*models.UserDeviceInfo, error)
	UpdateDeviceInfo(deviceID string, appVersion, osVersion, deviceModel string) error
	DeactivateDevice(deviceID string) error
	RemoveDevice(deviceID string) error
	UpdateDeviceLastUsed(deviceID string) error

	// Notification info methods
	GetUserNotificationInfo(userID string) (*models.UserNotificationInfo, error)
}
