package slack

import "context"

// SlackService interface defines methods for Slack notifications
type SlackService interface {
	SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
}
