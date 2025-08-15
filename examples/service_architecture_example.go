package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/services"
)

// Example demonstrating the new service architecture with interfaces

func main() {
	// Example 1: Basic service container usage
	basicExample()

	// Example 2: Configurable service container
	configurableExample()

	// Example 3: Mock service for testing
	mockExample()

	// Example 4: Service factory usage
	factoryExample()
}

// basicExample demonstrates basic service container usage
func basicExample() {
	fmt.Println("=== Basic Service Container Example ===")

	// Create service container
	container := services.NewServiceContainer()

	// Get services from container (returns interfaces)
	userService := container.GetUserService()
	_ = container.GetNotificationService() // Use notification service

	// Use services
	ctx := context.Background()

	// Create a user
	user := models.NewUser("john.doe@example.com", "John Doe")
	err := userService.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	fmt.Printf("Created user: %s\n", user.ID)

	// Get user by ID
	retrievedUser, err := userService.GetUserByID(ctx, user.ID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}

	fmt.Printf("Retrieved user: %s - %s\n", retrievedUser.ID, retrievedUser.Email)

	// Register a device
	device := &models.UserDeviceInfo{
		UserID:      user.ID,
		DeviceToken: "test_device_token_123",
		DeviceType:  "ios",
		AppVersion:  "1.0.0",
		OSVersion:   "iOS 16.0",
		DeviceModel: "iPhone 14",
	}

	createdDevice, err := userService.RegisterDevice(ctx, user.ID, device)
	if err != nil {
		log.Printf("Error registering device: %v", err)
		return
	}

	fmt.Printf("Registered device: %s\n", createdDevice.ID)

	// Get user notification info
	notificationInfo, err := userService.GetUserNotificationInfo(ctx, user.ID)
	if err != nil {
		log.Printf("Error getting notification info: %v", err)
		return
	}

	fmt.Printf("User has %d devices\n", len(notificationInfo.Devices))

	// Gracefully shutdown
	if err := container.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

// configurableExample demonstrates configurable service container usage
func configurableExample() {
	fmt.Println("\n=== Configurable Service Container Example ===")

	// Create custom configuration
	config := services.DefaultServiceConfig()
	config.EmailConfig.SMTPHost = "smtp.company.com"
	config.EmailConfig.SMTPPort = 587
	config.EmailConfig.FromEmail = "notifications@company.com"
	config.SlackConfig.BotToken = "xoxb-your-bot-token"
	config.SlackConfig.DefaultChannel = "#general"
	config.UserConfig.DatabaseURL = "postgres://localhost/notification_service"
	config.UserConfig.CacheTTL = 1800 // 30 minutes

	// Create container with configuration
	container := services.NewServiceContainerWithConfig(config)

	// Get services (now configured)
	_ = container.GetEmailService() // Use email service
	_ = container.GetSlackService() // Use slack service
	userService := container.GetUserService()

	fmt.Printf("Email service configured for: %s:%d\n", config.EmailConfig.SMTPHost, config.EmailConfig.SMTPPort)
	fmt.Printf("Slack service configured for channel: %s\n", config.SlackConfig.DefaultChannel)
	fmt.Printf("User service configured for database: %s\n", config.UserConfig.DatabaseURL)

	// Use configured services
	ctx := context.Background()

	// Create a test user
	user := models.NewUser("jane.smith@example.com", "Jane Smith")
	err := userService.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	fmt.Printf("Created configured user: %s\n", user.ID)

	// Gracefully shutdown
	if err := container.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

// MockUserService demonstrates how to create mock services for testing
type MockUserService struct {
	users   map[string]*models.User
	devices map[string]*models.UserDeviceInfo
}

// NewMockUserService creates a new mock user service
func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:   make(map[string]*models.User),
		devices: make(map[string]*models.UserDeviceInfo),
	}
}

// Implement UserService interface
func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID string) error {
	if _, exists := m.users[userID]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(m.users, userID)
	return nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserService) GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error) {
	var users []*models.User
	for _, id := range userIDs {
		if user, exists := m.users[id]; exists {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserService) RegisterDevice(ctx context.Context, userID string, device *models.UserDeviceInfo) (*models.UserDeviceInfo, error) {
	if _, exists := m.users[userID]; !exists {
		return nil, fmt.Errorf("user not found")
	}
	device.ID = fmt.Sprintf("device-%d", len(m.devices)+1)
	device.UserID = userID
	m.devices[device.ID] = device
	return device, nil
}

func (m *MockUserService) GetUserDevices(ctx context.Context, userID string) ([]*models.UserDeviceInfo, error) {
	if _, exists := m.users[userID]; !exists {
		return nil, fmt.Errorf("user not found")
	}
	var devices []*models.UserDeviceInfo
	for _, device := range m.devices {
		if device.UserID == userID {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (m *MockUserService) GetActiveUserDevices(ctx context.Context, userID string) ([]*models.UserDeviceInfo, error) {
	if _, exists := m.users[userID]; !exists {
		return nil, fmt.Errorf("user not found")
	}
	var devices []*models.UserDeviceInfo
	for _, device := range m.devices {
		if device.UserID == userID && device.IsActive {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (m *MockUserService) UpdateDeviceInfo(ctx context.Context, deviceID string, updates map[string]interface{}) (*models.UserDeviceInfo, error) {
	if device, exists := m.devices[deviceID]; exists {
		// Apply updates (simplified)
		if appVersion, ok := updates["app_version"].(string); ok {
			device.AppVersion = appVersion
		}
		if osVersion, ok := updates["os_version"].(string); ok {
			device.OSVersion = osVersion
		}
		return device, nil
	}
	return nil, fmt.Errorf("device not found")
}

func (m *MockUserService) DeactivateDevice(ctx context.Context, deviceID string) error {
	if device, exists := m.devices[deviceID]; exists {
		device.IsActive = false
		return nil
	}
	return fmt.Errorf("device not found")
}

func (m *MockUserService) RemoveDevice(ctx context.Context, deviceID string) error {
	if _, exists := m.devices[deviceID]; exists {
		delete(m.devices, deviceID)
		return nil
	}
	return fmt.Errorf("device not found")
}

func (m *MockUserService) UpdateDeviceLastUsed(ctx context.Context, deviceID string) error {
	if device, exists := m.devices[deviceID]; exists {
		device.UpdateLastUsed()
		return nil
	}
	return fmt.Errorf("device not found")
}

func (m *MockUserService) GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error) {
	user, err := m.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	devices, err := m.GetUserDevices(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.ToNotificationInfo(devices), nil
}

// mockExample demonstrates using mock services for testing
func mockExample() {
	fmt.Println("\n=== Mock Service Example ===")

	// Create mock service
	mockUserService := NewMockUserService()

	// Use mock service (same interface as real service)
	ctx := context.Background()

	// Create a user
	user := models.NewUser("test@example.com", "Test User")
	err := mockUserService.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	fmt.Printf("Created user in mock: %s\n", user.ID)

	// Register a device
	device := &models.UserDeviceInfo{
		UserID:      user.ID,
		DeviceToken: "mock_device_token",
		DeviceType:  "android",
		AppVersion:  "1.0.0",
		OSVersion:   "Android 13",
		DeviceModel: "Samsung Galaxy S23",
	}

	createdDevice, err := mockUserService.RegisterDevice(ctx, user.ID, device)
	if err != nil {
		log.Printf("Error registering device: %v", err)
		return
	}

	fmt.Printf("Registered device in mock: %s\n", createdDevice.ID)

	// Get user notification info
	notificationInfo, err := mockUserService.GetUserNotificationInfo(ctx, user.ID)
	if err != nil {
		log.Printf("Error getting notification info: %v", err)
		return
	}

	fmt.Printf("Mock user has %d devices\n", len(notificationInfo.Devices))
	fmt.Printf("Mock user email: %s\n", notificationInfo.Email)
}

// factoryExample demonstrates service factory usage
func factoryExample() {
	fmt.Println("\n=== Service Factory Example ===")

	// Create service factory
	factory := services.NewServiceFactory()

	// Create individual services
	userService := factory.NewUserService()
	emailService := factory.NewEmailService()
	slackService := factory.NewSlackService()
	inAppService := factory.NewInAppService()
	schedulerService := factory.NewSchedulerService()

	// Create notification manager with dependencies
	notificationService := factory.NewNotificationManager(
		emailService,
		slackService,
		inAppService,
		schedulerService,
	)

	fmt.Printf("Created %T\n", userService)
	fmt.Printf("Created %T\n", emailService)
	fmt.Printf("Created %T\n", slackService)
	fmt.Printf("Created %T\n", inAppService)
	fmt.Printf("Created %T\n", schedulerService)
	fmt.Printf("Created %T\n", notificationService)

	// Use services
	ctx := context.Background()

	// Create a user
	user := models.NewUser("factory@example.com", "Factory User")
	err := userService.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	fmt.Printf("Created user via factory: %s\n", user.ID)

	// Get all users
	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return
	}

	fmt.Printf("Total users via factory: %d\n", len(users))
}
