package email

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MockEmailServiceImpl implements the EmailService interface for testing/mock purposes
type MockEmailServiceImpl struct {
	outputPath string
}

// NewMockEmailService creates a new mock email service instance
func NewMockEmailService() EmailService {
	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	return &MockEmailServiceImpl{
		outputPath: filepath.Join(outputDir, "email.txt"),
	}
}

// SendEmail writes email notification to file instead of sending actual email
func (es *MockEmailServiceImpl) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*struct {
		ID       string
		Type     string
		Content  map[string]interface{}
		Template *struct {
			ID   string
			Data map[string]interface{}
		}
		Recipients []string
	})
	if !ok {
		return nil, ErrEmailSendFailed
	}

	// Create mock response
	response := &struct {
		ID      string    `json:"id"`
		Status  string    `json:"status"`
		Message string    `json:"message"`
		SentAt  time.Time `json:"sent_at"`
		Channel string    `json:"channel"`
	}{
		ID:      notif.ID,
		Status:  "mock_sent",
		Message: "Email notification written to file (mock mode)",
		SentAt:  time.Now(),
		Channel: "email",
	}

	// Prepare notification data for file output
	notificationData := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"id":         notif.ID,
		"type":       notif.Type,
		"content":    notif.Content,
		"recipients": notif.Recipients,
		"template":   notif.Template,
		"status":     "mock_sent",
		"channel":    "email",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(notificationData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification data: %w", err)
	}

	// Write to file
	file, err := os.OpenFile(es.outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Add separator and newline
	output := fmt.Sprintf("=== EMAIL NOTIFICATION ===\n%s\n\n", string(jsonData))
	if _, err := file.WriteString(output); err != nil {
		return nil, fmt.Errorf("failed to write to output file: %w", err)
	}

	return response, nil
}
