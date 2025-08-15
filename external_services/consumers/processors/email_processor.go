package consumers

import (
	"context"
	"fmt"
	"log"
)

// emailProcessor handles email notifications
type emailProcessor struct{}

// NewEmailProcessor creates a new email processor
func NewEmailProcessor() NotificationProcessor {
	return &emailProcessor{}
}

// ProcessNotification processes an email notification
func (ep *emailProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	fmt.Println("---------------> Processing email notification", message.Payload)
	log.Printf("---------------> Processing email notification: %s", message.Payload)
	// TODO: Implement actual email sending logic
	// This would integrate with your email service
	return nil
}

// GetNotificationType returns the notification type this processor handles
func (ep *emailProcessor) GetNotificationType() NotificationType {
	return EmailNotification
}
