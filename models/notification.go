package models

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	EmailNotification NotificationType = "email"
	SlackNotification NotificationType = "slack"
	InAppNotification NotificationType = "in_app"
)

// NotificationResponse represents the response after sending a notification
type NotificationResponse struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
	Channel string    `json:"channel"`
}
