package consumers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotificationProcessor is a mock implementation of NotificationProcessor
type MockNotificationProcessor struct {
	mock.Mock
}

func (m *MockNotificationProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockNotificationProcessor) GetNotificationType() NotificationType {
	args := m.Called()
	return args.Get(0).(NotificationType)
}

func TestWorkerProcessNotification(t *testing.T) {
	// Create a mock processor
	mockProcessor := new(MockNotificationProcessor)

	// Set up expectations
	mockProcessor.On("GetNotificationType").Return(EmailNotification)
	mockProcessor.On("ProcessNotification", mock.Anything, mock.Anything).Return(nil)

	// Create a channel for testing
	channel := make(chan string, 1)

	// Create a worker
	worker := NewWorker(channel, mockProcessor)

	// Start the worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := worker.Start(ctx)
	assert.NoError(t, err)

	// Send a test message
	testMessage := "test notification message"
	channel <- testMessage

	// Wait a bit for processing
	time.Sleep(100 * time.Millisecond)

	// Verify that ProcessNotification was called
	mockProcessor.AssertExpectations(t)

	// Stop the worker
	err = worker.Stop()
	assert.NoError(t, err)
}
