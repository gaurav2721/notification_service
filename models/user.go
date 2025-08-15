package models

import (
	"time"

	"github.com/google/uuid"
)

// UserDeviceInfo represents device information for InApp notifications
type UserDeviceInfo struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DeviceToken string    `json:"device_token"`
	DeviceType  string    `json:"device_type"` // "ios", "android", "web"
	AppVersion  string    `json:"app_version,omitempty"`
	OSVersion   string    `json:"os_version,omitempty"`
	DeviceModel string    `json:"device_model,omitempty"`
	IsActive    bool      `json:"is_active"`
	LastUsedAt  time.Time `json:"last_used_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User represents a user with essential information for notifications
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	SlackUserID  string    `json:"slack_user_id,omitempty"`
	SlackChannel string    `json:"slack_channel,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserNotificationInfo represents essential user info for notifications
type UserNotificationInfo struct {
	ID           string            `json:"id"`
	Email        string            `json:"email"`
	FullName     string            `json:"full_name"`
	SlackUserID  string            `json:"slack_user_id,omitempty"`
	SlackChannel string            `json:"slack_channel,omitempty"`
	PhoneNumber  string            `json:"phone_number,omitempty"`
	Devices      []*UserDeviceInfo `json:"devices,omitempty"`
}

// NewUser creates a new user with default values
func NewUser(email, fullName string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		FullName:  fullName,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToNotificationInfo converts User to UserNotificationInfo
func (u *User) ToNotificationInfo() *UserNotificationInfo {
	return &UserNotificationInfo{
		ID:           u.ID,
		Email:        u.Email,
		FullName:     u.FullName,
		SlackUserID:  u.SlackUserID,
		SlackChannel: u.SlackChannel,
		PhoneNumber:  u.PhoneNumber,
	}
}

// GetNotificationChannels returns enabled notification channels for the user
func (u *User) GetNotificationChannels() []string {
	var channels []string

	// Email is always available if user has email
	if u.Email != "" {
		channels = append(channels, "email")
	}

	// Slack is available if user has SlackUserID
	if u.SlackUserID != "" {
		channels = append(channels, "slack")
	}

	// InApp is always available (devices will be checked separately)
	channels = append(channels, "in_app")

	return channels
}

// NewUserDeviceInfo creates a new device info with default values
func NewUserDeviceInfo(userID, deviceToken, deviceType string) *UserDeviceInfo {
	now := time.Now()
	return &UserDeviceInfo{
		ID:          uuid.New().String(),
		UserID:      userID,
		DeviceToken: deviceToken,
		DeviceType:  deviceType,
		IsActive:    true,
		LastUsedAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateLastUsed updates the last used timestamp
func (d *UserDeviceInfo) UpdateLastUsed() {
	d.LastUsedAt = time.Now()
	d.UpdatedAt = time.Now()
}

// Deactivate marks the device as inactive
func (d *UserDeviceInfo) Deactivate() {
	d.IsActive = false
	d.UpdatedAt = time.Now()
}

// GetDeviceTokens returns all device tokens for a list of devices
func GetDeviceTokens(devices []*UserDeviceInfo) []string {
	var tokens []string
	for _, device := range devices {
		if device.IsActive && device.DeviceToken != "" {
			tokens = append(tokens, device.DeviceToken)
		}
	}
	return tokens
}

// GetDeviceTokensByType returns device tokens filtered by device type
func GetDeviceTokensByType(devices []*UserDeviceInfo, deviceType string) []string {
	var tokens []string
	for _, device := range devices {
		if device.IsActive && device.DeviceToken != "" && device.DeviceType == deviceType {
			tokens = append(tokens, device.DeviceToken)
		}
	}
	return tokens
}
