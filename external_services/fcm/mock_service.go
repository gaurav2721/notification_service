package fcm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gaurav2721/notification-service/models"
)

// MockFCMServiceImpl implements the FCMService interface for testing/mock purposes
type MockFCMServiceImpl struct {
	outputPath string
}

// NewMockFCMService creates a new mock FCM service instance
func NewMockFCMService() FCMService {
	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	return &MockFCMServiceImpl{
		outputPath: filepath.Join(outputDir, "fcm.txt"),
	}
}

// SendPushNotification writes FCM notification to file instead of sending actual push notification
func (fcm *MockFCMServiceImpl) SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*models.FCMNotificationRequest)
	if !ok {
		return nil, ErrInvalidNotificationPayload
	}

	// Validate the FCM notification
	if err := models.ValidateFCMNotification(notif); err != nil {
		return nil, fmt.Errorf("FCM validation failed: %w", err)
	}

	// Create mock response
	response := &models.FCMResponse{
		ID:           notif.ID,
		Status:       "mock_sent",
		Message:      "FCM notification written to file (mock mode)",
		SentAt:       time.Now(),
		Channel:      "fcm",
		SuccessCount: 1,
		FailureCount: 0,
	}

	// Prepare notification data for file output
	notificationData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"id":        notif.ID,
		"type":      notif.Type,
		"content":   notif.Content,
		"recipient": notif.Recipient,
		"status":    "mock_sent",
		"channel":   "fcm",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(notificationData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification data: %w", err)
	}

	// Write to file
	file, err := os.OpenFile(fcm.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Add separator and newline
	output := fmt.Sprintf("=== FCM NOTIFICATION ===\n%s\n\n", string(jsonData))
	if _, err := file.WriteString(output); err != nil {
		return nil, fmt.Errorf("failed to write to output file: %w", err)
	}

	return response, nil
}
