package consumers

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// workerPool represents a pool of workers for a specific notification type
type workerPool struct {
	notificationType NotificationType
	channel          chan string
	processor        NotificationProcessor
	workers          []ConsumerWorker
	workerCount      int
	running          bool
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	mu               sync.RWMutex
}

// NewWorkerPool creates a new worker pool for a specific notification type
func NewWorkerPool(
	notificationType NotificationType,
	channel chan string,
	processor NotificationProcessor,
	workerCount int,
) ConsumerWorkerPool {
	return &workerPool{
		notificationType: notificationType,
		channel:          channel,
		processor:        processor,
		workerCount:      workerCount,
		workers:          make([]ConsumerWorker, 0, workerCount),
		running:          false,
	}
}

// Start initializes and starts the worker pool
func (wp *workerPool) Start(ctx context.Context) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return fmt.Errorf("worker pool for %s is already running", wp.notificationType)
	}

	wp.ctx, wp.cancel = context.WithCancel(ctx)
	wp.running = true

	// Create and start workers
	for i := 0; i < wp.workerCount; i++ {
		worker := NewWorker(wp.channel, wp.processor)
		wp.workers = append(wp.workers, worker)

		wp.wg.Add(1)
		go func(w ConsumerWorker) {
			defer wp.wg.Done()
			if err := w.Start(wp.ctx); err != nil {
				log.Printf("Failed to start worker %s: %v", w.GetWorkerID(), err)
			}
		}(worker)
	}

	log.Printf("Started worker pool for %s with %d workers", wp.notificationType, wp.workerCount)
	return nil
}

// Stop gracefully shuts down the worker pool
func (wp *workerPool) Stop() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.running {
		return nil
	}

	wp.running = false
	if wp.cancel != nil {
		wp.cancel()
	}

	// Stop all workers
	for _, worker := range wp.workers {
		if err := worker.Stop(); err != nil {
			log.Printf("Error stopping worker %s: %v", worker.GetWorkerID(), err)
		}
	}

	// Wait for all workers to finish
	wp.wg.Wait()

	log.Printf("Stopped worker pool for %s", wp.notificationType)
	return nil
}

// GetWorkerCount returns the current number of active workers
func (wp *workerPool) GetWorkerCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return len(wp.workers)
}

// GetNotificationType returns the type of notifications this pool handles
func (wp *workerPool) GetNotificationType() NotificationType {
	return wp.notificationType
}

// IsRunning returns true if the worker pool is currently running
func (wp *workerPool) IsRunning() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.running
}

// GetChannel returns the channel this pool reads from
func (wp *workerPool) GetChannel() chan string {
	return wp.channel
}

// UpdateWorkerCount updates the number of workers in the pool
func (wp *workerPool) UpdateWorkerCount(count int) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return fmt.Errorf("cannot update worker count while pool is running")
	}

	wp.workerCount = count
	wp.workers = make([]ConsumerWorker, 0, count)
	return nil
}
