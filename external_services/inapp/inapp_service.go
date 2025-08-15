package inapp

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// InAppServiceImpl implements the InAppService interface
type InAppServiceImpl struct {
	notifications map[string][]interface{}
	mutex         sync.RWMutex
}

// NewInAppService creates a new in-app notification service instance
func NewInAppService() InAppService {
	return &InAppServiceImpl{
		notifications: make(map[string][]interface{}),
	}
}

// NewInAppServiceWithConfig creates a new in-app service with custom configuration
func NewInAppServiceWithConfig(config *InAppConfig) InAppService {
	return &InAppServiceImpl{
		notifications: make(map[string][]interface{}),
	}
}

// SendInAppNotification sends an in-app notification
func (ias *InAppServiceImpl) SendInAppNotification(ctx context.Context, notification interface{}) (interface{}, error) {
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
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrInAppSendFailed
	}

	// Store notification for each recipient
	ias.mutex.Lock()
	defer ias.mutex.Unlock()

	for _, recipient := range notif.Recipients {
		if ias.notifications[recipient] == nil {
			ias.notifications[recipient] = make([]interface{}, 0)
		}
		ias.notifications[recipient] = append(ias.notifications[recipient], notif)
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
		Message: "In-app notification sent successfully",
		SentAt:  time.Now(),
		Channel: "in_app",
	}, nil
}

// GetNotificationsForUser retrieves notifications for a specific user
func (ias *InAppServiceImpl) GetNotificationsForUser(userID string) []interface{} {
	ias.mutex.RLock()
	defer ias.mutex.RUnlock()

	if notifications, exists := ias.notifications[userID]; exists {
		return notifications
	}
	return []interface{}{}
}

// MarkNotificationAsRead marks a notification as read for a user
func (ias *InAppServiceImpl) MarkNotificationAsRead(userID, notificationID string) error {
	ias.mutex.Lock()
	defer ias.mutex.Unlock()

	if notifications, exists := ias.notifications[userID]; exists {
		for i, notification := range notifications {
			if notif, ok := notification.(*struct {
				ID       string
				Type     string
				Content  map[string]interface{}
				Template *struct {
					ID   string
					Data map[string]interface{}
				}
				Recipients  []string
				Metadata    map[string]interface{}
				ScheduledAt *time.Time
			}); ok && notif.ID == notificationID {
				// Remove the notification from the user's list
				ias.notifications[userID] = append(notifications[:i], notifications[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("notification not found for user")
}
