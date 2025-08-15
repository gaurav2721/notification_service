package consumers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
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

	log.Printf("Worker %s started for %s notifications", w.id, w.processor.GetNotificationType())
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

	log.Printf("Worker %s stopped", w.id)
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

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %s received shutdown signal", w.id)
			return

		case message, ok := <-w.channel:
			if !ok {
				log.Printf("Worker %s: channel closed", w.id)
				return
			}

			fmt.Println("---------> gaurav message", message)

			// Process the notification
			if err := w.processMessage(message); err != nil {
				log.Printf("Worker %s: error processing message: %v", w.id, err)
				// Continue processing other messages even if one fails
			}
		}
	}
}

// processMessage processes a single notification message
func (w *worker) processMessage(message string) error {
	start := time.Now()

	// Parse the message into NotificationMessage
	// This is a simplified version - in a real implementation,
	// you might want to use JSON unmarshaling or a more robust parsing mechanism
	notificationMsg := NotificationMessage{
		Type:      w.processor.GetNotificationType(),
		Payload:   message,
		ID:        uuid.New().String(),
		Timestamp: time.Now().Unix(),
	}

	// Process the notification using the processor
	if err := w.processor.ProcessNotification(w.ctx, notificationMsg); err != nil {
		return fmt.Errorf("failed to process notification: %w", err)
	}

	log.Printf("Worker %s processed notification in %v", w.id, time.Since(start))
	return nil
}
