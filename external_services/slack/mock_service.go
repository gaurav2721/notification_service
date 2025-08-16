package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gaurav2721/notification-service/models"
)

// MockSlackServiceImpl implements the SlackService interface for testing/mock purposes
type MockSlackServiceImpl struct {
	outputPath string
}

// NewMockSlackService creates a new mock Slack service instance
func NewMockSlackService() SlackService {
	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	return &MockSlackServiceImpl{
		outputPath: filepath.Join(outputDir, "slack.txt"),
	}
}

// SendSlackMessage writes Slack notification to file instead of sending actual message
func (ss *MockSlackServiceImpl) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*models.SlackNotificationRequest)
	if !ok {
		return nil, ErrSlackSendFailed
	}

	// Validate the slack notification
	if err := models.ValidateSlackNotification(notif); err != nil {
		return nil, fmt.Errorf("slack validation failed: %w", err)
	}

	// Create mock response
	response := &models.SlackResponse{
		ID:      notif.ID,
		Status:  "mock_sent",
		Message: "Slack notification written to file (mock mode)",
		SentAt:  time.Now(),
		Channel: "slack",
	}

	// Prepare notification data for file output
	notificationData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"id":        notif.ID,
		"type":      notif.Type,
		"content":   notif.Content,
		"recipient": notif.Recipient,
		"status":    "mock_sent",
		"channel":   "slack",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(notificationData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification data: %w", err)
	}

	// Write to file
	file, err := os.OpenFile(ss.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Add separator and newline
	output := fmt.Sprintf("=== SLACK NOTIFICATION ===\n%s\n\n", string(jsonData))
	if _, err := file.WriteString(output); err != nil {
		return nil, fmt.Errorf("failed to write to output file: %w", err)
	}

	return response, nil
}
