package apns

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gaurav2721/notification-service/models"
)

// MockAPNSServiceImpl implements the APNSService interface for testing/mock purposes
type MockAPNSServiceImpl struct {
	outputPath string
}

// NewMockAPNSService creates a new mock APNS service instance
func NewMockAPNSService() APNSService {
	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	return &MockAPNSServiceImpl{
		outputPath: filepath.Join(outputDir, "apns.txt"),
	}
}

// SendPushNotification writes APNS notification to file instead of sending actual push notification
func (aps *MockAPNSServiceImpl) SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*models.APNSNotificationRequest)
	if !ok {
		return nil, ErrInvalidNotificationPayload
	}

	// Validate the APNS notification
	if err := models.ValidateAPNSNotification(notif); err != nil {
		return nil, fmt.Errorf("APNS validation failed: %w", err)
	}

	// Create mock response
	response := &models.APNSResponse{
		ID:           notif.ID,
		Status:       "mock_sent",
		Message:      "APNS notification written to file (mock mode)",
		SentAt:       time.Now(),
		Channel:      "apns",
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
		"channel":   "apns",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(notificationData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification data: %w", err)
	}

	// Write to file
	file, err := os.OpenFile(aps.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Add separator and newline
	output := fmt.Sprintf("=== APNS NOTIFICATION ===\n%s\n\n", string(jsonData))
	if _, err := file.WriteString(output); err != nil {
		return nil, fmt.Errorf("failed to write to output file: %w", err)
	}

	return response, nil
}
