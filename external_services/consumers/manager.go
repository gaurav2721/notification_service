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

	logrus.WithField("worker_pools", len(cm.workerPools)).Info("Consumer manager initialized successfully")
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

	logrus.Info("Consumer manager started all worker pools")
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

	logrus.Info("Consumer manager stopped all worker pools")
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
	processor := NewEmailProcessor()
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
	processor := NewSlackProcessor()
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
	processor := NewIOSPushProcessor()
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
	processor := NewAndroidPushProcessor()
	pool := NewWorkerPool(
		AndroidPushNotification,
		cm.config.KafkaService.GetAndroidPushNotificationChannel(),
		processor,
		cm.config.AndroidPushWorkerCount,
	)
	cm.workerPools[AndroidPushNotification] = pool
}

// Processor factory functions - these reference the implementations in the processors package
// NewEmailProcessor creates a new email processor
func NewEmailProcessor() NotificationProcessor {
	return &emailProcessor{}
}

// NewSlackProcessor creates a new slack processor
func NewSlackProcessor() NotificationProcessor {
	return &slackProcessor{}
}

// NewIOSPushProcessor creates a new iOS push notification processor
func NewIOSPushProcessor() NotificationProcessor {
	return &iosPushProcessor{}
}

// NewAndroidPushProcessor creates a new Android push notification processor
func NewAndroidPushProcessor() NotificationProcessor {
	return &androidPushProcessor{}
}

// Processor implementations - these are the actual implementations with debug logs
type emailProcessor struct{}

func (ep *emailProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
	}).Info("-----------------> gaurav singh Processing email notification")

	// TODO: Implement actual email sending logic
	// This would integrate with your email service
	return nil
}

func (ep *emailProcessor) GetNotificationType() NotificationType {
	return EmailNotification
}

type slackProcessor struct{}

func (sp *slackProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
	}).Info("Processing slack notification")

	// TODO: Implement actual slack message sending logic
	// This would integrate with your slack service
	return nil
}

func (sp *slackProcessor) GetNotificationType() NotificationType {
	return SlackNotification
}

type iosPushProcessor struct{}

func (ip *iosPushProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
	}).Info("Processing iOS push notification")

	// TODO: Implement actual iOS push notification logic
	// This would integrate with your iOS push service
	return nil
}

func (ip *iosPushProcessor) GetNotificationType() NotificationType {
	return IOSPushNotification
}

type androidPushProcessor struct{}

func (ap *androidPushProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	logrus.WithFields(logrus.Fields{
		"notification_id": message.ID,
		"type":            message.Type,
		"payload":         message.Payload,
	}).Info("Processing Android push notification")

	// TODO: Implement actual Android push notification logic
	// This would integrate with your Android push service
	return nil
}

func (ap *androidPushProcessor) GetNotificationType() NotificationType {
	return AndroidPushNotification
}
