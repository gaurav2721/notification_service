package user

import (
	"testing"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	service := NewUserService()
	assert.NotNil(t, service)
}

func TestUserService_GetUserByID(t *testing.T) {
	service := NewUserService()

	// Test getting existing user
	user, err := service.GetUserByID("user-001")
	require.NoError(t, err)
	assert.Equal(t, "user-001", user.ID)
	assert.Equal(t, "john.doe@company.com", user.Email)
	assert.Equal(t, "John Doe", user.FullName)

	// Test getting non-existent user
	_, err = service.GetUserByID("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserService_GetUsersByIDs(t *testing.T) {
	service := NewUserService()

	// Test getting multiple existing users
	userIDs := []string{"user-001", "user-002", "user-003"}
	users, err := service.GetUsersByIDs(userIDs)
	require.NoError(t, err)
	assert.Len(t, users, 3)

	// Test getting mix of existing and non-existing users
	userIDs = []string{"user-001", "non-existent", "user-002"}
	users, err = service.GetUsersByIDs(userIDs)
	require.NoError(t, err)
	assert.Len(t, users, 2) // Should only return existing users
}

func TestUserService_GetAllUsers(t *testing.T) {
	service := NewUserService()

	// Get all users
	users, err := service.GetAllUsers()
	require.NoError(t, err)
	assert.Len(t, users, 8) // All 8 preloaded users are active

	// Verify all returned users are active
	for _, user := range users {
		assert.True(t, user.IsActive)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	service := NewUserService()

	// Create a new user
	newUser := &models.User{
		ID:       "user-new-001",
		Email:    "newuser@company.com",
		FullName: "New User",
		IsActive: true,
	}

	err := service.CreateUser(newUser)
	require.NoError(t, err)

	// Verify user was created
	createdUser, err := service.GetUserByID("user-new-001")
	require.NoError(t, err)
	assert.Equal(t, "newuser@company.com", createdUser.Email)
	assert.Equal(t, "New User", createdUser.FullName)
}

func TestUserService_UpdateUser(t *testing.T) {
	service := NewUserService()

	// Get existing user
	user, err := service.GetUserByID("user-001")
	require.NoError(t, err)

	// Update user
	user.FullName = "John Doe Updated"
	err = service.UpdateUser(user)
	require.NoError(t, err)

	// Verify update
	updatedUser, err := service.GetUserByID("user-001")
	require.NoError(t, err)
	assert.Equal(t, "John Doe Updated", updatedUser.FullName)
}

func TestUserService_DeleteUser(t *testing.T) {
	service := NewUserService()

	// Delete user
	err := service.DeleteUser("user-001")
	require.NoError(t, err)

	// Verify user is inactive
	_, err = service.GetUserByID("user-001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user is inactive")
}

func TestUserService_RegisterDevice(t *testing.T) {
	service := NewUserService()

	// Register a new device
	device, err := service.RegisterDevice("user-001", "test_device_token", "ios")
	require.NoError(t, err)
	assert.Equal(t, "user-001", device.UserID)
	assert.Equal(t, "test_device_token", device.DeviceToken)
	assert.Equal(t, "ios", device.DeviceType)
	assert.True(t, device.IsActive)

	// Test registering same device again (should update last used)
	device2, err := service.RegisterDevice("user-001", "test_device_token", "ios")
	require.NoError(t, err)
	assert.Equal(t, device.ID, device2.ID)
}

func TestUserService_GetUserDevices(t *testing.T) {
	service := NewUserService()

	// Get all devices for user
	devices, err := service.GetUserDevices("user-001")
	require.NoError(t, err)
	assert.Len(t, devices, 2) // user-001 has 2 devices in sample data
}

func TestUserService_GetActiveUserDevices(t *testing.T) {
	service := NewUserService()

	// Get active devices for user
	devices, err := service.GetActiveUserDevices("user-001")
	require.NoError(t, err)
	assert.Len(t, devices, 2) // user-001 has 2 active devices

	// Get active devices for user with inactive device
	devices, err = service.GetActiveUserDevices("user-003")
	require.NoError(t, err)
	assert.Len(t, devices, 0) // user-003 has 1 inactive device
}

func TestUserService_UpdateDeviceInfo(t *testing.T) {
	service := NewUserService()

	// Get a device first
	devices, err := service.GetUserDevices("user-001")
	require.NoError(t, err)
	require.Len(t, devices, 2)

	deviceID := devices[0].ID

	// Update device info
	err = service.UpdateDeviceInfo(deviceID, "1.3.0", "iOS 17.0", "iPhone 15")
	require.NoError(t, err)

	// Verify update
	devices, err = service.GetUserDevices("user-001")
	require.NoError(t, err)

	var updatedDevice *models.UserDeviceInfo
	for _, device := range devices {
		if device.ID == deviceID {
			updatedDevice = device
			break
		}
	}
	require.NotNil(t, updatedDevice)
	assert.Equal(t, "1.3.0", updatedDevice.AppVersion)
	assert.Equal(t, "iOS 17.0", updatedDevice.OSVersion)
	assert.Equal(t, "iPhone 15", updatedDevice.DeviceModel)
}

func TestUserService_DeactivateDevice(t *testing.T) {
	service := NewUserService()

	// Get a device first
	devices, err := service.GetUserDevices("user-001")
	require.NoError(t, err)
	require.Len(t, devices, 2)

	deviceID := devices[0].ID

	// Deactivate device
	err = service.DeactivateDevice(deviceID)
	require.NoError(t, err)

	// Verify device is inactive
	devices, err = service.GetActiveUserDevices("user-001")
	require.NoError(t, err)
	assert.Len(t, devices, 1) // Should have 1 active device now
}

func TestUserService_RemoveDevice(t *testing.T) {
	service := NewUserService()

	// Get a device first
	devices, err := service.GetUserDevices("user-001")
	require.NoError(t, err)
	require.Len(t, devices, 2)

	deviceID := devices[0].ID

	// Remove device
	err = service.RemoveDevice(deviceID)
	require.NoError(t, err)

	// Verify device is removed
	devices, err = service.GetUserDevices("user-001")
	require.NoError(t, err)
	assert.Len(t, devices, 1) // Should have 1 device now
}

func TestUserService_UpdateDeviceLastUsed(t *testing.T) {
	service := NewUserService()

	// Get a device first
	devices, err := service.GetUserDevices("user-001")
	require.NoError(t, err)
	require.Len(t, devices, 2)

	deviceID := devices[0].ID
	originalLastUsed := devices[0].LastUsedAt

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Update last used
	err = service.UpdateDeviceLastUsed(deviceID)
	require.NoError(t, err)

	// Verify last used was updated
	devices, err = service.GetUserDevices("user-001")
	require.NoError(t, err)

	var updatedDevice *models.UserDeviceInfo
	for _, device := range devices {
		if device.ID == deviceID {
			updatedDevice = device
			break
		}
	}
	require.NotNil(t, updatedDevice)
	assert.True(t, updatedDevice.LastUsedAt.After(originalLastUsed))
}

func TestUserService_GetUserNotificationInfo(t *testing.T) {
	service := NewUserService()

	// Get notification info
	info, err := service.GetUserNotificationInfo("user-001")
	require.NoError(t, err)
	assert.Equal(t, "user-001", info.ID)
	assert.Equal(t, "john.doe@company.com", info.Email)
	assert.Equal(t, "John Doe", info.FullName)
	assert.Len(t, info.Devices, 2) // Should include devices
}

func TestUser_GetNotificationChannels(t *testing.T) {
	// Test user with all channels available
	user := &models.User{
		Email:       "test@company.com",
		SlackUserID: "U1234567890",
	}

	channels := user.GetNotificationChannels()
	assert.Len(t, channels, 3)
	assert.Contains(t, channels, "email")
	assert.Contains(t, channels, "slack")
	assert.Contains(t, channels, "in_app")

	// Test user with no slack user ID
	user.SlackUserID = ""
	channels = user.GetNotificationChannels()
	assert.Len(t, channels, 2)
	assert.Contains(t, channels, "email")
	assert.Contains(t, channels, "in_app")
	assert.NotContains(t, channels, "slack")

	// Test user with no email
	user.Email = ""
	channels = user.GetNotificationChannels()
	assert.Len(t, channels, 1)
	assert.Contains(t, channels, "in_app")
	assert.NotContains(t, channels, "email")
	assert.NotContains(t, channels, "slack")
}

func TestUser_ToNotificationInfo(t *testing.T) {
	user := &models.User{
		ID:           "user-001",
		Email:        "test@company.com",
		FullName:     "Test User",
		SlackUserID:  "U1234567890",
		SlackChannel: "#test",
		PhoneNumber:  "+1-555-0101",
	}

	info := user.ToNotificationInfo()
	assert.Equal(t, user.ID, info.ID)
	assert.Equal(t, user.Email, info.Email)
	assert.Equal(t, user.FullName, info.FullName)
	assert.Equal(t, user.SlackUserID, info.SlackUserID)
	assert.Equal(t, user.SlackChannel, info.SlackChannel)
	assert.Equal(t, user.PhoneNumber, info.PhoneNumber)
}

func TestUserDeviceInfo_HelperFunctions(t *testing.T) {
	devices := []*models.UserDeviceInfo{
		{
			ID:          "device-1",
			DeviceToken: "token-1",
			DeviceType:  "ios",
			IsActive:    true,
		},
		{
			ID:          "device-2",
			DeviceToken: "token-2",
			DeviceType:  "android",
			IsActive:    true,
		},
		{
			ID:          "device-3",
			DeviceToken: "token-3",
			DeviceType:  "ios",
			IsActive:    false, // inactive
		},
		{
			ID:          "device-4",
			DeviceToken: "", // empty token
			DeviceType:  "web",
			IsActive:    true,
		},
	}

	// Test GetDeviceTokens
	tokens := models.GetDeviceTokens(devices)
	assert.Len(t, tokens, 2) // Only active devices with tokens
	assert.Contains(t, tokens, "token-1")
	assert.Contains(t, tokens, "token-2")

	// Test GetDeviceTokensByType
	iosTokens := models.GetDeviceTokensByType(devices, "ios")
	assert.Len(t, iosTokens, 1) // Only active iOS device with token
	assert.Contains(t, iosTokens, "token-1")

	androidTokens := models.GetDeviceTokensByType(devices, "android")
	assert.Len(t, androidTokens, 1)
	assert.Contains(t, androidTokens, "token-2")

	webTokens := models.GetDeviceTokensByType(devices, "web")
	assert.Len(t, webTokens, 0) // No active web device with token
}
