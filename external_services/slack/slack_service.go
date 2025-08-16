package slack

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/slack-go/slack"
)

// SlackServiceImpl implements the SlackService interface
type SlackServiceImpl struct {
	client  *slack.Client
	channel string
}

// NewSlackService creates a new Slack service instance
// It checks environment variables and returns mock service if config is incomplete
func NewSlackService() SlackService {
	token := os.Getenv("SLACK_BOT_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL_ID")

	// Check if all required environment variables are present and non-empty
	if token == "" || channel == "" {
		return NewMockSlackService()
	}

	client := slack.New(token)

	return &SlackServiceImpl{
		client:  client,
		channel: channel,
	}
}

// SendSlackMessage sends a Slack notification
func (ss *SlackServiceImpl) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*models.SlackNotificationRequest)
	if !ok {
		return nil, ErrSlackSendFailed
	}

	// Validate the slack notification
	if err := models.ValidateSlackNotification(notif); err != nil {
		return nil, fmt.Errorf("slack validation failed: %w", err)
	}

	// Extract text from content
	text := notif.Content.Text

	// Create Slack message
	msg := slack.MsgOptionText(text, false)

	// Send message
	_, _, err := ss.client.PostMessage(ss.channel, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send slack message: %w", err)
	}

	// Return success response
	return &models.SlackResponse{
		ID:      notif.ID,
		Status:  "sent",
		Message: "Slack message sent successfully",
		SentAt:  time.Now(),
		Channel: "slack",
	}, nil
}
