package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/services"
)

// ExampleNotificationService demonstrates integration with UserService
type ExampleNotificationService struct {
	userService  services.UserService
	emailService services.EmailService
	slackService services.SlackService
	inAppService services.InAppService
}

// NewExampleNotificationService creates a new notification service with user integration
func NewExampleNotificationService() *ExampleNotificationService {
	return &ExampleNotificationService{
		userService:  services.NewUserService(),
		emailService: &MockEmailService{},
		slackService: &MockSlackService{},
		inAppService: &MockInAppService{},
	}
}

// SendNotificationToUsers sends notifications to multiple users
func (s *ExampleNotificationService) SendNotificationToUsers(ctx context.Context, userIDs []string, notification *models.Notification) error {
	// Get user information for all users
	userInfos, err := s.userService.GetUsersNotificationInfo(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	fmt.Printf("Sending notification to %d users\n", len(userInfos))

	for _, userInfo := range userInfos {
		// Get user's notification channels
		channels := userInfo.GetNotificationChannels()

		fmt.Printf("Sending to user %s via channels: %v\n", userInfo.FullName, channels)

		// Send to each available channel
		for _, channel := range channels {
			switch channel {
			case "email":
				if err := s.sendEmail(userInfo, notification); err != nil {
					log.Printf("Failed to send email to %s: %v", userInfo.Email, err)
				}
			case "slack":
				if err := s.sendSlack(userInfo, notification); err != nil {
					log.Printf("Failed to send slack to %s: %v", userInfo.SlackUserID, err)
				}
			case "in_app":
				if err := s.sendInApp(userInfo, notification); err != nil {
					log.Printf("Failed to send in-app to %s: %v", userInfo.ID, err)
				}
			}
		}
	}

	return nil
}

// sendEmail sends email notification to a user
func (s *ExampleNotificationService) sendEmail(userInfo *models.UserNotificationInfo, notification *models.Notification) error {
	// Create email-specific notification
	emailNotification := &EmailNotification{
		To:      userInfo.Email,
		Subject: notification.Title,
		Body:    notification.Message,
		User:    userInfo,
	}

	_, err := s.emailService.SendEmail(context.Background(), emailNotification)
	return err
}

// sendSlack sends slack notification to a user
func (s *ExampleNotificationService) sendSlack(userInfo *models.UserNotificationInfo, notification *models.Notification) error {
	// Create slack-specific notification
	slackNotification := &SlackNotification{
		UserID:  userInfo.SlackUserID,
		Channel: userInfo.SlackChannel,
		Message: notification.Message,
		User:    userInfo,
	}

	_, err := s.slackService.SendSlackMessage(context.Background(), slackNotification)
	return err
}

// sendInApp sends in-app notification to a user
func (s *ExampleNotificationService) sendInApp(userInfo *models.UserNotificationInfo, notification *models.Notification) error {
	// Create in-app notification
	inAppNotification := &InAppNotification{
		UserID:  userInfo.ID,
		Title:   notification.Title,
		Message: notification.Message,
		User:    userInfo,
	}

	_, err := s.inAppService.SendInAppNotification(context.Background(), inAppNotification)
	return err
}

// Example usage functions
func ExampleUsage() {
	ctx := context.Background()
	notificationService := NewExampleNotificationService()

	// Example 1: Send notification to specific users
	fmt.Println("=== Example 1: Send to specific users ===")
	notification := models.NewNotification(
		models.EmailNotification,
		"Project Update",
		"Your project has been updated with new features.",
		[]string{"user-001", "user-002", "user-003"},
	)

	err := notificationService.SendNotificationToUsers(ctx, []string{"user-001", "user-002"}, notification)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	// Example 2: Send urgent notification to specific users
	fmt.Println("\n=== Example 2: Send urgent notification ===")
	urgentNotification := models.NewNotification(
		models.SlackNotification,
		"URGENT: System Maintenance",
		"System maintenance will begin in 30 minutes.",
		[]string{},
	)
	urgentNotification.Priority = models.UrgentPriority

	err = notificationService.SendNotificationToUsers(ctx, []string{"user-001", "user-005"}, urgentNotification)
	if err != nil {
		log.Printf("Error sending urgent notification: %v", err)
	}

	// Example 3: Demonstrate device management
	fmt.Println("\n=== Example 3: Device management ===")

	// Register a new device
	device, err := notificationService.userService.RegisterDevice(ctx, "user-001", "new_device_token", "android")
	if err != nil {
		log.Printf("Error registering device: %v", err)
	} else {
		fmt.Printf("Registered device: %s for user: %s\n", device.ID, device.UserID)
	}

	// Get user's devices
	devices, err := notificationService.userService.GetActiveUserDevices(ctx, "user-001")
	if err != nil {
		log.Printf("Error getting devices: %v", err)
	} else {
		fmt.Printf("User user-001 has %d active devices\n", len(devices))
		for _, device := range devices {
			fmt.Printf("  - %s (%s): %s\n", device.DeviceType, device.DeviceModel, device.DeviceToken)
		}
	}

	// Example 4: Get user notification info
	fmt.Println("\n=== Example 4: User notification info ===")
	userInfo, err := notificationService.userService.GetUserNotificationInfo(ctx, "user-001")
	if err != nil {
		log.Printf("Error getting user info: %v", err)
	} else {
		fmt.Printf("User: %s (%s)\n", userInfo.FullName, userInfo.Email)
		fmt.Printf("Available channels: %v\n", userInfo.GetNotificationChannels())
		fmt.Printf("Active devices: %d\n", len(userInfo.Devices))
	}
}

// Mock services for demonstration
type MockEmailService struct{}
type MockSlackService struct{}
type MockInAppService struct{}

type EmailNotification struct {
	To      string
	Subject string
	Body    string
	User    *models.UserNotificationInfo
}

type SlackNotification struct {
	UserID  string
	Channel string
	Message string
	User    *models.UserNotificationInfo
}

type InAppNotification struct {
	UserID  string
	Title   string
	Message string
	User    *models.UserNotificationInfo
}

func (m *MockEmailService) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	if emailNotif, ok := notification.(*EmailNotification); ok {
		fmt.Printf("📧 Email sent to %s: %s\n", emailNotif.To, emailNotif.Subject)
	}
	return "email_sent", nil
}

func (m *MockSlackService) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	if slackNotif, ok := notification.(*SlackNotification); ok {
		fmt.Printf("💬 Slack message sent to %s in %s: %s\n", slackNotif.UserID, slackNotif.Channel, slackNotif.Message)
	}
	return "slack_sent", nil
}

func (m *MockInAppService) SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	if inAppNotif, ok := notification.(*InAppNotification); ok {
		fmt.Printf("🔔 In-app notification sent to %s: %s\n", inAppNotif.UserID, inAppNotif.Title)
	}
	return "inapp_sent", nil
}

func main() {
	ExampleUsage()
}
