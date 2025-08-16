package apns

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
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
	notif, ok := notification.(*struct {
		ID       string
		Type     string
		Content  map[string]interface{}
		Template *struct {
			ID   string
			Data map[string]interface{}
		}
		Recipients  []string
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrInvalidNotificationPayload
	}

	// Create mock response
	response := &struct {
		ID           string    `json:"id"`
		Status       string    `json:"status"`
		Message      string    `json:"message"`
		SentAt       time.Time `json:"sent_at"`
		Channel      string    `json:"channel"`
		SuccessCount int       `json:"success_count"`
		FailureCount int       `json:"failure_count"`
	}{
		ID:           notif.ID,
		Status:       "mock_sent",
		Message:      "APNS notification written to file (mock mode)",
		SentAt:       time.Now(),
		Channel:      "apns",
		SuccessCount: len(notif.Recipients),
		FailureCount: 0,
	}

	// Prepare notification data for file output
	notificationData := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"id":           notif.ID,
		"type":         notif.Type,
		"content":      notif.Content,
		"recipients":   notif.Recipients,
		"template":     notif.Template,
		"scheduled_at": notif.ScheduledAt,
		"status":       "mock_sent",
		"channel":      "apns",
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
