package consumers

import (
	"context"
	"fmt"
	"log"
	"sync"
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
		return fmt.Errorf("consumer manager is already initialized")
	}

	// Create worker pools for each notification type
	cm.createEmailWorkerPool()
	cm.createSlackWorkerPool()
	cm.createIOSPushWorkerPool()
	cm.createAndroidPushWorkerPool()

	log.Printf("Consumer manager initialized with %d worker pools", len(cm.workerPools))
	return nil
}

// Start starts all consumer worker pools
func (cm *consumerManager) Start(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.running {
		return fmt.Errorf("consumer manager is already running")
	}

	cm.ctx, cm.cancel = context.WithCancel(ctx)
	cm.running = true

	// Start all worker pools
	for notificationType, pool := range cm.workerPools {
		cm.wg.Add(1)
		go func(nt NotificationType, p ConsumerWorkerPool) {
			defer cm.wg.Done()
			if err := p.Start(cm.ctx); err != nil {
				log.Printf("Failed to start worker pool for %s: %v", nt, err)
			}
		}(notificationType, pool)
	}

	log.Printf("Consumer manager started all worker pools")
	return nil
}

// Stop gracefully shuts down all consumer worker pools
func (cm *consumerManager) Stop() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.running {
		return nil
	}

	cm.running = false
	if cm.cancel != nil {
		cm.cancel()
	}

	// Stop all worker pools
	for notificationType, pool := range cm.workerPools {
		if err := pool.Stop(); err != nil {
			log.Printf("Error stopping worker pool for %s: %v", notificationType, err)
		}
	}

	// Wait for all worker pools to finish
	cm.wg.Wait()

	log.Printf("Consumer manager stopped all worker pools")
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
	processor := NewEmailProcessor() // This would be implemented separately
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
	processor := NewSlackProcessor() // This would be implemented separately
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
	processor := NewIOSPushProcessor() // This would be implemented separately
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
	processor := NewAndroidPushProcessor() // This would be implemented separately
	pool := NewWorkerPool(
		AndroidPushNotification,
		cm.config.KafkaService.GetAndroidPushNotificationChannel(),
		processor,
		cm.config.AndroidPushWorkerCount,
	)
	cm.workerPools[AndroidPushNotification] = pool
}
