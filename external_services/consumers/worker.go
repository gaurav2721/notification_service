package consumers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// worker represents a single worker that processes notifications
type worker struct {
	id        string
	channel   chan string
	processor NotificationProcessor
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.RWMutex
}

// NewWorker creates a new worker instance
func NewWorker(channel chan string, processor NotificationProcessor) ConsumerWorker {
	return &worker{
		id:        uuid.New().String(),
		channel:   channel,
		processor: processor,
		running:   false,
	}
}

// Start begins processing notifications from the channel
func (w *worker) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return fmt.Errorf("worker %s is already running", w.id)
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.running = true
	w.wg.Add(1)

	go w.processLoop()

	logrus.WithFields(logrus.Fields{
		"worker_id": w.id,
		"type":      w.processor.GetNotificationType(),
	}).Info("Worker started")
	return nil
}

// Stop gracefully shuts down the worker
func (w *worker) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return nil
	}

	w.running = false
	if w.cancel != nil {
		w.cancel()
	}

	// Wait for the worker to finish processing
	w.wg.Wait()

	logrus.WithField("worker_id", w.id).Info("Worker stopped")
	return nil
}

// IsRunning returns true if the worker is currently running
func (w *worker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

// GetWorkerID returns the unique identifier for this worker
func (w *worker) GetWorkerID() string {
	return w.id
}

// processLoop is the main processing loop for the worker
func (w *worker) processLoop() {
	defer w.wg.Done()

	logrus.WithFields(logrus.Fields{
		"worker_id": w.id,
		"type":      w.processor.GetNotificationType(),
	}).Debug("Worker processing loop started")

	for {
		select {
		case <-w.ctx.Done():
			logrus.WithField("worker_id", w.id).Debug("Worker received shutdown signal")
			return

		case message, ok := <-w.channel:
			if !ok {
				logrus.WithField("worker_id", w.id).Debug("Worker: channel closed")
				return
			}

			logrus.WithFields(logrus.Fields{
				"worker_id": w.id,
				"message":   message,
			}).Debug("Worker received message from channel")

			// Process the notification
			if err := w.processMessage(message); err != nil {
				logrus.WithFields(logrus.Fields{
					"worker_id": w.id,
					"error":     err.Error(),
				}).Error("Worker error processing message")
				// Continue processing other messages even if one fails
			}
		}
	}
}

// processMessage processes a single notification message
func (w *worker) processMessage(message string) error {
	start := time.Now()

	logrus.WithFields(logrus.Fields{
		"worker_id": w.id,
		"message":   message,
	}).Debug("Worker starting to process message")

	// Parse the message into NotificationMessage
	// This is a simplified version - in a real implementation,
	// you might want to use JSON unmarshaling or a more robust parsing mechanism
	notificationMsg := NotificationMessage{
		Type:      w.processor.GetNotificationType(),
		Payload:   message,
		ID:        uuid.New().String(),
		Timestamp: time.Now().Unix(),
	}

	logrus.WithFields(logrus.Fields{
		"worker_id":       w.id,
		"notification_id": notificationMsg.ID,
		"type":            notificationMsg.Type,
		"payload":         notificationMsg.Payload,
	}).Debug("Worker created notification message")

	// Process the notification using the processor
	logrus.WithFields(logrus.Fields{
		"worker_id":       w.id,
		"notification_id": notificationMsg.ID,
		"type":            notificationMsg.Type,
	}).Debug("Worker calling ProcessNotification")

	if err := w.processor.ProcessNotification(w.ctx, notificationMsg); err != nil {
		logrus.WithFields(logrus.Fields{
			"worker_id":       w.id,
			"notification_id": notificationMsg.ID,
			"error":           err.Error(),
		}).Error("Worker ProcessNotification failed")
		return fmt.Errorf("failed to process notification: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"worker_id":       w.id,
		"notification_id": notificationMsg.ID,
		"duration":        time.Since(start),
	}).Info("Worker processed notification successfully")
	return nil
}
