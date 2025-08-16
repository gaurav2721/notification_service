package user

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gaurav2721/notification-service/models"
)

// UserService interface defines methods for user management
type UserService interface {
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error

	// Device management methods
	RegisterDevice(ctx context.Context, userID, deviceToken, deviceType string) (*models.UserDeviceInfo, error)
	GetUserDevices(ctx context.Context, userID string) ([]*models.UserDeviceInfo, error)
	GetActiveUserDevices(ctx context.Context, userID string) ([]*models.UserDeviceInfo, error)
	UpdateDeviceInfo(ctx context.Context, deviceID string, appVersion, osVersion, deviceModel string) error
	DeactivateDevice(ctx context.Context, deviceID string) error
	RemoveDevice(ctx context.Context, deviceID string) error
	UpdateDeviceLastUsed(ctx context.Context, deviceID string) error

	// Notification info methods
	GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error)
}

// userService implements UserService interface
type userService struct {
	users   map[string]*models.User
	devices map[string]*models.UserDeviceInfo // deviceID -> UserDeviceInfo
	mutex   sync.RWMutex
}

// NewUserService creates a new user service with preloaded data
func NewUserService() UserService {
	service := &userService{
		users:   make(map[string]*models.User),
		devices: make(map[string]*models.UserDeviceInfo),
	}

	// Preload users with sample data
	service.preloadUsers()
	// Preload devices with sample data
	service.preloadDevices()

	return service
}

// preloadUsers loads sample user data for testing and development
func (s *userService) preloadUsers() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Sample users with essential notification information
	users := []*models.User{
		{
			ID:           "user-001",
			Email:        "john.doe@company.com",
			FullName:     "John Doe",
			SlackUserID:  "U1234567890",
			SlackChannel: "#general",
			PhoneNumber:  "+1-555-0101",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -6, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-002",
			Email:        "jane.smith@company.com",
			FullName:     "Jane Smith",
			SlackUserID:  "U0987654321",
			SlackChannel: "#design",
			PhoneNumber:  "+1-555-0102",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -4, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-003",
			Email:        "mike.johnson@company.com",
			FullName:     "Mike Johnson",
			SlackUserID:  "U1122334455",
			SlackChannel: "#marketing",
			PhoneNumber:  "+1-555-0103",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -8, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-004",
			Email:        "sarah.wilson@company.com",
			FullName:     "Sarah Wilson",
			SlackUserID:  "U5566778899",
			SlackChannel: "#sales",
			PhoneNumber:  "+1-555-0104",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -2, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-005",
			Email:        "david.brown@company.com",
			FullName:     "David Brown",
			SlackUserID:  "U9988776655",
			SlackChannel: "#engineering",
			PhoneNumber:  "+1-555-0105",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -12, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-006",
			Email:        "lisa.garcia@company.com",
			FullName:     "Lisa Garcia",
			SlackUserID:  "U4433221100",
			SlackChannel: "#marketing",
			PhoneNumber:  "+1-555-0106",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -10, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-007",
			Email:        "robert.taylor@company.com",
			FullName:     "Robert Taylor",
			SlackUserID:  "U1122334455",
			SlackChannel: "#sales",
			PhoneNumber:  "+1-555-0107",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -9, 0),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "user-008",
			Email:        "emma.davis@company.com",
			FullName:     "Emma Davis",
			SlackUserID:  "U6677889900",
			SlackChannel: "#executives",
			PhoneNumber:  "+1-555-0108",
			IsActive:     true,
			CreatedAt:    time.Now().AddDate(0, -15, 0),
			UpdatedAt:    time.Now(),
		},
	}

	// Add users to the map
	for _, user := range users {
		s.users[user.ID] = user
	}
}

// preloadDevices loads sample device data for testing and development
func (s *userService) preloadDevices() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Sample devices for InApp notifications
	devices := []*models.UserDeviceInfo{
		{
			ID:          "device-001",
			UserID:      "user-001",
			DeviceToken: "ios_token_123456789",
			DeviceType:  "ios",
			AppVersion:  "1.2.3",
			OSVersion:   "iOS 16.0",
			DeviceModel: "iPhone 14",
			IsActive:    true,
			LastUsedAt:  time.Now(),
			CreatedAt:   time.Now().AddDate(0, -2, 0),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "device-002",
			UserID:      "user-001",
			DeviceToken: "android_token_987654321",
			DeviceType:  "android",
			AppVersion:  "1.2.3",
			OSVersion:   "Android 13",
			DeviceModel: "Samsung Galaxy S23",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Hour),
			CreatedAt:   time.Now().AddDate(0, -1, 0),
			UpdatedAt:   time.Now().Add(-time.Hour),
		},
		{
			ID:          "device-003",
			UserID:      "user-002",
			DeviceToken: "web_token_456789123",
			DeviceType:  "web",
			AppVersion:  "1.2.3",
			OSVersion:   "Chrome 120.0",
			DeviceModel: "Desktop",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Minute * 30),
			CreatedAt:   time.Now().AddDate(0, -3, 0),
			UpdatedAt:   time.Now().Add(-time.Minute * 30),
		},
		{
			ID:          "device-004",
			UserID:      "user-003",
			DeviceToken: "ios_token_789123456",
			DeviceType:  "ios",
			AppVersion:  "1.2.2",
			OSVersion:   "iOS 15.5",
			DeviceModel: "iPhone 13",
			IsActive:    false, // Inactive device
			LastUsedAt:  time.Now().AddDate(0, 0, -7),
			CreatedAt:   time.Now().AddDate(0, -6, 0),
			UpdatedAt:   time.Now().AddDate(0, 0, -7),
		},
		{
			ID:          "device-005",
			UserID:      "user-004",
			DeviceToken: "android_token_555666777",
			DeviceType:  "android",
			AppVersion:  "1.2.4",
			OSVersion:   "Android 14",
			DeviceModel: "Google Pixel 8",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Minute * 15),
			CreatedAt:   time.Now().AddDate(0, -1, 0),
			UpdatedAt:   time.Now().Add(-time.Minute * 15),
		},
		{
			ID:          "device-006",
			UserID:      "user-005",
			DeviceToken: "ios_token_111222333",
			DeviceType:  "ios",
			AppVersion:  "1.2.3",
			OSVersion:   "iOS 17.0",
			DeviceModel: "iPhone 15 Pro",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Hour * 2),
			CreatedAt:   time.Now().AddDate(0, -3, 0),
			UpdatedAt:   time.Now().Add(-time.Hour * 2),
		},
		{
			ID:          "device-007",
			UserID:      "user-006",
			DeviceToken: "web_token_444555666",
			DeviceType:  "web",
			AppVersion:  "1.2.3",
			OSVersion:   "Firefox 121.0",
			DeviceModel: "MacBook Pro",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Minute * 45),
			CreatedAt:   time.Now().AddDate(0, -2, 0),
			UpdatedAt:   time.Now().Add(-time.Minute * 45),
		},
		{
			ID:          "device-008",
			UserID:      "user-007",
			DeviceToken: "android_token_777888999",
			DeviceType:  "android",
			AppVersion:  "1.2.1",
			OSVersion:   "Android 12",
			DeviceModel: "OnePlus 9",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Hour * 3),
			CreatedAt:   time.Now().AddDate(0, -4, 0),
			UpdatedAt:   time.Now().Add(-time.Hour * 3),
		},
		{
			ID:          "device-009",
			UserID:      "user-008",
			DeviceToken: "ios_token_000111222",
			DeviceType:  "ios",
			AppVersion:  "1.2.5",
			OSVersion:   "iOS 17.2",
			DeviceModel: "iPad Pro",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Minute * 10),
			CreatedAt:   time.Now().AddDate(0, -1, 0),
			UpdatedAt:   time.Now().Add(-time.Minute * 10),
		},
		{
			ID:          "device-010",
			UserID:      "user-002",
			DeviceToken: "android_token_333444555",
			DeviceType:  "android",
			AppVersion:  "1.2.3",
			OSVersion:   "Android 13",
			DeviceModel: "Samsung Galaxy Tab S9",
			IsActive:    true,
			LastUsedAt:  time.Now().Add(-time.Hour * 4),
			CreatedAt:   time.Now().AddDate(0, -2, 0),
			UpdatedAt:   time.Now().Add(-time.Hour * 4),
		},
	}

	// Add devices to the map
	for _, device := range devices {
		s.devices[device.ID] = device
	}
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	return user, nil
}

// GetUsersByIDs retrieves multiple users by their IDs
func (s *userService) GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var users []*models.User
	for _, userID := range userIDs {
		if user, exists := s.users[userID]; exists && user.IsActive {
			users = append(users, user)
		}
	}

	return users, nil
}

// GetAllUsers retrieves all users from the service
func (s *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var allUsers []*models.User
	for _, user := range s.users {
		if user.IsActive {
			allUsers = append(allUsers, user)
		}
	}

	return allUsers, nil
}

// CreateUser adds a new user to the service
func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user

	return nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	s.users[user.ID] = user

	return nil
}

// DeleteUser removes a user from the service (soft delete)
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if user, exists := s.users[userID]; exists {
		user.IsActive = false
		user.UpdatedAt = time.Now()
		return nil
	}

	return errors.New("user not found")
}

// RegisterDevice registers a new device for a user
func (s *userService) RegisterDevice(ctx context.Context, userID, deviceToken, deviceType string) (*models.UserDeviceInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate user exists and is active
	if _, exists := s.users[userID]; !exists {
		return nil, ErrUserNotFound
	}

	// Device token validation - only check if it's not empty
	if deviceToken == "" {
		return nil, errors.New("device token cannot be empty")
	}

	// Check if device already exists for this user
	for _, device := range s.devices {
		if device.UserID == userID && device.DeviceToken == deviceToken {
			// Update existing device
			device.UpdateLastUsed()
			return device, nil
		}
	}

	// Create new device
	device := models.NewUserDeviceInfo(userID, deviceToken, deviceType)
	s.devices[device.ID] = device

	return device, nil
}

// GetUserDevices retrieves all devices for a user
func (s *userService) GetUserDevices(ctx context.Context, userID string) ([]*models.UserDeviceInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var devices []*models.UserDeviceInfo
	for _, device := range s.devices {
		if device.UserID == userID {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

// GetActiveUserDevices retrieves all active devices for a user
func (s *userService) GetActiveUserDevices(ctx context.Context, userID string) ([]*models.UserDeviceInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var devices []*models.UserDeviceInfo
	for _, device := range s.devices {
		if device.UserID == userID && device.IsActive {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

// UpdateDeviceInfo updates device information
func (s *userService) UpdateDeviceInfo(ctx context.Context, deviceID string, appVersion, osVersion, deviceModel string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device, exists := s.devices[deviceID]
	if !exists {
		return ErrDeviceNotFound
	}

	if appVersion != "" {
		device.AppVersion = appVersion
	}
	if osVersion != "" {
		device.OSVersion = osVersion
	}
	if deviceModel != "" {
		device.DeviceModel = deviceModel
	}

	device.UpdateLastUsed()
	return nil
}

// DeactivateDevice marks a device as inactive
func (s *userService) DeactivateDevice(ctx context.Context, deviceID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device, exists := s.devices[deviceID]
	if !exists {
		return ErrDeviceNotFound
	}

	device.Deactivate()
	return nil
}

// RemoveDevice completely removes a device
func (s *userService) RemoveDevice(ctx context.Context, deviceID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.devices[deviceID]; !exists {
		return ErrDeviceNotFound
	}

	delete(s.devices, deviceID)
	return nil
}

// UpdateDeviceLastUsed updates the last used timestamp for a device
func (s *userService) UpdateDeviceLastUsed(ctx context.Context, deviceID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device, exists := s.devices[deviceID]
	if !exists {
		return ErrDeviceNotFound
	}

	device.UpdateLastUsed()
	return nil
}

// GetUserNotificationInfo retrieves essential user info for notifications
func (s *userService) GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get active devices for the user
	devices, err := s.GetActiveUserDevices(ctx, userID)
	if err != nil {
		return nil, err
	}

	notificationInfo := user.ToNotificationInfo()
	notificationInfo.Devices = devices

	return notificationInfo, nil
}
