package notification_manager

import (
	"sync"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/sirupsen/logrus"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusScheduled NotificationStatus = "scheduled"
	StatusQueued    NotificationStatus = "queued"
	StatusSent      NotificationStatus = "sent"
	StatusFailed    NotificationStatus = "failed"
	StatusCancelled NotificationStatus = "cancelled"
)

// NotificationRecord represents a stored notification record
type NotificationRecord struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
	Template    *models.TemplateData   `json:"template,omitempty"`
	Recipients  []string               `json:"recipients"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	From        *struct {
		Email string `json:"email"`
	} `json:"from,omitempty"`
	Status    NotificationStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
	Error     string             `json:"error,omitempty"`
}

// InMemoryStorage provides thread-safe in-memory storage for notifications
type InMemoryStorage struct {
	notifications map[string]*NotificationRecord
	mutex         sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		notifications: make(map[string]*NotificationRecord),
	}
}

// StoreNotification stores a notification record
func (s *InMemoryStorage) StoreNotification(notificationID string, notification *models.NotificationRequest) error {
	if notificationID == "" {
		return ErrUnsupportedNotificationType
	}

	if notification == nil {
		return ErrUnsupportedNotificationType
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	record := &NotificationRecord{
		ID:          notificationID,
		Type:        notification.Type,
		Content:     notification.Content,
		Template:    notification.Template,
		Recipients:  notification.Recipients,
		ScheduledAt: notification.ScheduledAt,
		From:        notification.From,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.notifications[notificationID] = record

	logrus.WithFields(logrus.Fields{
		"notification_id": notificationID,
		"type":            notification.Type,
		"status":          StatusPending,
	}).Debug("Notification stored in memory")

	return nil
}

// GetNotification retrieves a notification record by ID
func (s *InMemoryStorage) GetNotification(notificationID string) (*NotificationRecord, error) {
	if notificationID == "" {
		return nil, ErrUnsupportedNotificationType
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	record, exists := s.notifications[notificationID]
	if !exists {
		return nil, ErrUnsupportedNotificationType
	}

	return record, nil
}

// UpdateNotificationStatus updates the status of a notification
func (s *InMemoryStorage) UpdateNotificationStatus(notificationID string, status NotificationStatus, errorMsg string) error {
	if notificationID == "" {
		return ErrUnsupportedNotificationType
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	record, exists := s.notifications[notificationID]
	if !exists {
		return ErrUnsupportedNotificationType
	}

	oldStatus := record.Status
	record.Status = status
	record.UpdatedAt = time.Now()
	record.Error = errorMsg

	// Set SentAt timestamp if status is sent
	if status == StatusSent {
		now := time.Now()
		record.SentAt = &now
	}

	logrus.WithFields(logrus.Fields{
		"notification_id": notificationID,
		"old_status":      oldStatus,
		"new_status":      status,
		"error":           errorMsg,
	}).Info("Notification status updated in memory")

	return nil
}

// GetAllNotifications retrieves all stored notifications
func (s *InMemoryStorage) GetAllNotifications() []*NotificationRecord {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	notifications := make([]*NotificationRecord, 0, len(s.notifications))
	for _, record := range s.notifications {
		notifications = append(notifications, record)
	}

	return notifications
}

// GetNotificationsByStatus retrieves notifications by status
func (s *InMemoryStorage) GetNotificationsByStatus(status NotificationStatus) []*NotificationRecord {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var notifications []*NotificationRecord
	for _, record := range s.notifications {
		if record.Status == status {
			notifications = append(notifications, record)
		}
	}

	return notifications
}

// DeleteNotification removes a notification from storage
func (s *InMemoryStorage) DeleteNotification(notificationID string) error {
	if notificationID == "" {
		return ErrUnsupportedNotificationType
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.notifications[notificationID]; !exists {
		return ErrUnsupportedNotificationType
	}

	delete(s.notifications, notificationID)

	logrus.WithField("notification_id", notificationID).Debug("Notification deleted from memory")

	return nil
}

// GetStorageStats returns basic statistics about the storage
func (s *InMemoryStorage) GetStorageStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_notifications"] = len(s.notifications)

	statusCounts := make(map[NotificationStatus]int)
	for _, record := range s.notifications {
		statusCounts[record.Status]++
	}
	stats["status_counts"] = statusCounts

	return stats
}
