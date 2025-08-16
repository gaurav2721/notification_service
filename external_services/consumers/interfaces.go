package consumers

import (
	"context"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	EmailNotification       NotificationType = "email"
	SlackNotification       NotificationType = "slack"
	IOSPushNotification     NotificationType = "ios_push"
	AndroidPushNotification NotificationType = "android_push"
)

// NotificationMessage represents a notification message from Kafka
type NotificationMessage struct {
	Type      NotificationType `json:"type"`
	Payload   string           `json:"payload"`
	ID        string           `json:"id"`
	Timestamp int64            `json:"timestamp"`
}

// ConsumerWorker represents a single worker that processes notifications
type ConsumerWorker interface {
	Start(ctx context.Context) error
	Stop() error
	IsRunning() bool
	GetWorkerID() string
}

// ConsumerWorkerPool represents a pool of workers for a specific notification type
type ConsumerWorkerPool interface {
	// Start initializes and starts the worker pool
	Start(ctx context.Context) error

	// Stop gracefully shuts down the worker pool
	Stop() error

	// GetWorkerCount returns the current number of active workers
	GetWorkerCount() int

	// GetNotificationType returns the type of notifications this pool handles
	GetNotificationType() NotificationType

	// IsRunning returns true if the worker pool is currently running
	IsRunning() bool

	// GetChannel returns the channel this pool reads from
	GetChannel() chan string
}

// ConsumerManager manages all consumer worker pools
type ConsumerManager interface {
	// Initialize creates and configures all consumer worker pools
	Initialize(ctx context.Context) error

	// Start starts all consumer worker pools
	Start(ctx context.Context) error

	// Stop gracefully shuts down all consumer worker pools
	Stop() error

	// GetWorkerPool returns a specific worker pool by notification type
	GetWorkerPool(notificationType NotificationType) (ConsumerWorkerPool, error)

	// GetAllWorkerPools returns all worker pools
	GetAllWorkerPools() map[NotificationType]ConsumerWorkerPool

	// GetStatus returns the status of all worker pools
	GetStatus() map[NotificationType]bool

	// UpdateWorkerCount updates the number of workers for a specific pool
	UpdateWorkerCount(notificationType NotificationType, count int) error
}

// ConsumerConfig holds configuration for consumer worker pools
type ConsumerConfig struct {
	EmailWorkerCount       int `json:"email_worker_count" env:"EMAIL_WORKER_COUNT" env-default:"5"`
	SlackWorkerCount       int `json:"slack_worker_count" env:"SLACK_WORKER_COUNT" env-default:"3"`
	IOSPushWorkerCount     int `json:"ios_push_worker_count" env:"IOS_PUSH_WORKER_COUNT" env-default:"3"`
	AndroidPushWorkerCount int `json:"android_push_worker_count" env:"ANDROID_PUSH_WORKER_COUNT" env-default:"3"`

	// Service dependencies
	EmailService interface {
		SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
	}
	SlackService interface {
		SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
	}
	APNSService interface {
		SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	}
	FCMService interface {
		SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	}

	// Kafka service interface for getting channels
	KafkaService interface {
		GetEmailChannel() chan string
		GetSlackChannel() chan string
		GetIOSPushNotificationChannel() chan string
		GetAndroidPushNotificationChannel() chan string
		Close()
	}
}

// NotificationProcessor defines the interface for processing notifications
type NotificationProcessor interface {
	ProcessNotification(ctx context.Context, message NotificationMessage) error
	GetNotificationType() NotificationType
}
