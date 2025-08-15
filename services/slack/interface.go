package slack

import (
	"context"
	"errors"
)

// SlackService interface defines methods for Slack notifications
type SlackService interface {
	SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
}

// SlackConfig holds Slack service configuration
type SlackConfig struct {
	BotToken       string
	DefaultChannel string
	WebhookURL     string
}

// DefaultSlackConfig returns default Slack configuration
func DefaultSlackConfig() *SlackConfig {
	return &SlackConfig{
		BotToken:       "",
		DefaultChannel: "#general",
		WebhookURL:     "",
	}
}

// Slack service errors
var (
	ErrSlackSendFailed   = errors.New("failed to send slack message")
	ErrInvalidChannel    = errors.New("invalid slack channel")
	ErrSlackTokenMissing = errors.New("slack bot token is missing")
)
