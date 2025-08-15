package slack

import "errors"

// Slack service errors
var (
	ErrSlackSendFailed   = errors.New("failed to send slack message")
	ErrInvalidChannel    = errors.New("invalid slack channel")
	ErrSlackTokenMissing = errors.New("slack bot token is missing")
)
