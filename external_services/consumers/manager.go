package consumers

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// consumerManager manages all consumer worker pools
type consumerManager struct {
	config      ConsumerConfig
	workerPools map[NotificationType]ConsumerWorkerPool
	running     bool
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

// NewConsumerManager creates a new consumer manager
func NewConsumerManager(config ConsumerConfig) ConsumerManager {
	return &consumerManager{
		config:      config,
		workerPools: make(map[NotificationType]ConsumerWorkerPool),
		running:     false,
	}
}

// NewConsumerManagerWithServices creates a new consumer manager with service dependencies
func NewConsumerManagerWithServices(
	emailService interface {
		SendEmail(ctx context.Context, notification interface{}) (interface{}, error)
	},
	slackService interface {
		SendSlackMessage(ctx context.Context, notification interface{}) (interface{}, error)
	},
	apnsService interface {
		SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	},
	fcmService interface {
		SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)
	},
	kafkaService interface {
		GetEmailChannel() chan string
		GetSlackChannel() chan string
		GetIOSPushNotificationChannel() chan string
		GetAndroidPushNotificationChannel() chan string
		Close()
	},
	config ConsumerConfig,
) ConsumerManager {
	// Update config with service dependencies
	config.EmailService = emailService
	config.SlackService = slackService
	config.APNSService = apnsService
	config.FCMService = fcmService
	config.KafkaService = kafkaService

	return &consumerManager{
		config:      config,
		workerPools: make(map[NotificationType]ConsumerWorkerPool),
		running:     false,
	}
}

// Initialize creates and configures all consumer worker pools
func (cm *consumerManager) Initialize(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.running {
		logrus.Warn("Consumer manager is already initialized")
		return fmt.Errorf("consumer manager is already initialized")
	}

	logrus.Debug("Initializing consumer manager")

	// Create worker pools for each notification type
	logrus.Debug("Creating email worker pool")
	cm.createEmailWorkerPool()

	logrus.Debug("Creating slack worker pool")
	cm.createSlackWorkerPool()

	logrus.Debug("Creating iOS push worker pool")
	cm.createIOSPushWorkerPool()

	logrus.Debug("Creating Android push worker pool")
	cm.createAndroidPushWorkerPool()

	logrus.WithField("worker_pools", len(cm.workerPools)).Debug("Consumer manager initialized successfully")
	return nil
}

// Start starts all consumer worker pools
func (cm *consumerManager) Start(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.running {
		logrus.Warn("Consumer manager is already running")
		return fmt.Errorf("consumer manager is already running")
	}

	logrus.Debug("Starting consumer manager")
	cm.ctx, cm.cancel = context.WithCancel(ctx)
	cm.running = true

	// Start all worker pools
	for notificationType, pool := range cm.workerPools {
		cm.wg.Add(1)
		go func(nt NotificationType, p ConsumerWorkerPool) {
			defer cm.wg.Done()
			logrus.WithField("notification_type", nt).Debug("Starting worker pool")
			if err := p.Start(cm.ctx); err != nil {
				logrus.WithFields(logrus.Fields{
					"notification_type": nt,
					"error":             err.Error(),
				}).Error("Failed to start worker pool")
			} else {
				logrus.WithField("notification_type", nt).Debug("Worker pool started successfully")
			}
		}(notificationType, pool)
	}

	logrus.Debug("Consumer manager started all worker pools")
	return nil
}

// Stop gracefully shuts down all consumer worker pools
func (cm *consumerManager) Stop() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.running {
		logrus.Debug("Consumer manager is not running")
		return nil
	}

	logrus.Debug("Stopping consumer manager")
	cm.running = false
	if cm.cancel != nil {
		cm.cancel()
	}

	// Stop all worker pools
	for notificationType, pool := range cm.workerPools {
		logrus.WithField("notification_type", notificationType).Debug("Stopping worker pool")
		if err := pool.Stop(); err != nil {
			logrus.WithFields(logrus.Fields{
				"notification_type": notificationType,
				"error":             err.Error(),
			}).Error("Error stopping worker pool")
		} else {
			logrus.WithField("notification_type", notificationType).Debug("Worker pool stopped successfully")
		}
	}

	// Wait for all worker pools to finish
	logrus.Debug("Waiting for all worker pools to finish")
	cm.wg.Wait()

	logrus.Debug("Consumer manager stopped all worker pools")
	return nil
}

// GetWorkerPool returns a specific worker pool by notification type
func (cm *consumerManager) GetWorkerPool(notificationType NotificationType) (ConsumerWorkerPool, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	pool, exists := cm.workerPools[notificationType]
	if !exists {
		return nil, fmt.Errorf("worker pool for %s not found", notificationType)
	}

	return pool, nil
}

// GetAllWorkerPools returns all worker pools
func (cm *consumerManager) GetAllWorkerPools() map[NotificationType]ConsumerWorkerPool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Create a copy to avoid race conditions
	pools := make(map[NotificationType]ConsumerWorkerPool)
	for k, v := range cm.workerPools {
		pools[k] = v
	}

	return pools
}

// GetStatus returns the status of all worker pools
func (cm *consumerManager) GetStatus() map[NotificationType]bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	status := make(map[NotificationType]bool)
	for notificationType, pool := range cm.workerPools {
		status[notificationType] = pool.IsRunning()
	}

	return status
}

// UpdateWorkerCount updates the number of workers for a specific pool
func (cm *consumerManager) UpdateWorkerCount(notificationType NotificationType, count int) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	pool, exists := cm.workerPools[notificationType]
	if !exists {
		return fmt.Errorf("worker pool for %s not found", notificationType)
	}

	// Cast to concrete type to access UpdateWorkerCount method
	if wp, ok := pool.(*workerPool); ok {
		return wp.UpdateWorkerCount(count)
	}

	return fmt.Errorf("worker pool for %s does not support updating worker count", notificationType)
}

// createEmailWorkerPool creates the email worker pool
func (cm *consumerManager) createEmailWorkerPool() {
	var processor NotificationProcessor

	// Use injected email service if available, otherwise create default
	if cm.config.EmailService != nil {
		processor = NewEmailProcessorWithService(cm.config.EmailService)
	} else {
		processor = NewEmailProcessor()
	}

	pool := NewWorkerPool(
		EmailNotification,
		cm.config.KafkaService.GetEmailChannel(),
		processor,
		cm.config.EmailWorkerCount,
	)
	cm.workerPools[EmailNotification] = pool
}

// createSlackWorkerPool creates the slack worker pool
func (cm *consumerManager) createSlackWorkerPool() {
	var processor NotificationProcessor

	// Use injected slack service if available, otherwise create default
	if cm.config.SlackService != nil {
		processor = NewSlackProcessorWithService(cm.config.SlackService)
	} else {
		processor = NewSlackProcessor()
	}

	pool := NewWorkerPool(
		SlackNotification,
		cm.config.KafkaService.GetSlackChannel(),
		processor,
		cm.config.SlackWorkerCount,
	)
	cm.workerPools[SlackNotification] = pool
}

// createIOSPushWorkerPool creates the iOS push notification worker pool
func (cm *consumerManager) createIOSPushWorkerPool() {
	var processor NotificationProcessor

	// Use injected APNS service if available, otherwise create default
	if cm.config.APNSService != nil {
		processor = NewIOSPushProcessorWithService(cm.config.APNSService)
	} else {
		processor = NewIOSPushProcessor()
	}

	pool := NewWorkerPool(
		IOSPushNotification,
		cm.config.KafkaService.GetIOSPushNotificationChannel(),
		processor,
		cm.config.IOSPushWorkerCount,
	)
	cm.workerPools[IOSPushNotification] = pool
}

// createAndroidPushWorkerPool creates the Android push notification worker pool
func (cm *consumerManager) createAndroidPushWorkerPool() {
	var processor NotificationProcessor

	// Use injected FCM service if available, otherwise create default
	if cm.config.FCMService != nil {
		processor = NewAndroidPushProcessorWithService(cm.config.FCMService)
	} else {
		processor = NewAndroidPushProcessor()
	}

	pool := NewWorkerPool(
		AndroidPushNotification,
		cm.config.KafkaService.GetAndroidPushNotificationChannel(),
		processor,
		cm.config.AndroidPushWorkerCount,
	)
	cm.workerPools[AndroidPushNotification] = pool
}
