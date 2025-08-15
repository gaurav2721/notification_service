package slack

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
)

// SlackServiceImpl implements the SlackService interface
type SlackServiceImpl struct {
	client  *slack.Client
	channel string
}

// NewSlackService creates a new Slack service instance
func NewSlackService() SlackService {
	token := os.Getenv("SLACK_BOT_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL_ID")

	client := slack.New(token)

	return &SlackServiceImpl{
		client:  client,
		channel: channel,
	}
}

// NewSlackServiceWithConfig creates a new Slack service with custom configuration
func NewSlackServiceWithConfig(config *SlackConfig) SlackService {
	client := slack.New(config.BotToken)

	return &SlackServiceImpl{
		client:  client,
		channel: config.DefaultChannel,
	}
}

// SendSlackMessage sends a Slack notification
func (ss *SlackServiceImpl) SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*struct {
		ID         string
		Type       string
		Title      string
		Message    string
		Recipients []string
		Metadata   map[string]interface{}
	})
	if !ok {
		return nil, ErrSlackSendFailed
	}

	// Create Slack message
	msg := slack.MsgOptionText(notif.Message, false)

	// Send message
	_, _, err := ss.client.PostMessage(ss.channel, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send slack message: %w", err)
	}

	// Return success response
	return &struct {
		ID      string    `json:"id"`
		Status  string    `json:"status"`
		Message string    `json:"message"`
		SentAt  time.Time `json:"sent_at"`
		Channel string    `json:"channel"`
	}{
		ID:      notif.ID,
		Status:  "sent",
		Message: "Slack message sent successfully",
		SentAt:  time.Now(),
		Channel: "slack",
	}, nil
}
