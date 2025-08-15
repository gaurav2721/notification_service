package kafka

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

// kafkaServiceImpl implements the KafkaService interface
type kafkaServiceImpl struct {
	emailChannel                   chan string
	slackChannel                   chan string
	iosPushNotificationChannel     chan string
	androidPushNotificationChannel chan string
	mu                             sync.RWMutex
	closed                         bool
}

// NewKafkaService creates a new instance of KafkaService
func NewKafkaService() (KafkaService, error) {
	// Read buffer sizes from environment variables with defaults
	emailBufferSize := getEnvAsInt("EMAIL_CHANNEL_BUFFER_SIZE", 100)
	slackBufferSize := getEnvAsInt("SLACK_CHANNEL_BUFFER_SIZE", 100)
	iosPushBufferSize := getEnvAsInt("IOS_PUSH_CHANNEL_BUFFER_SIZE", 100)
	androidPushBufferSize := getEnvAsInt("ANDROID_PUSH_CHANNEL_BUFFER_SIZE", 100)

	service := &kafkaServiceImpl{
		emailChannel:                   make(chan string, emailBufferSize),
		slackChannel:                   make(chan string, slackBufferSize),
		iosPushNotificationChannel:     make(chan string, iosPushBufferSize),
		androidPushNotificationChannel: make(chan string, androidPushBufferSize),
		closed:                         false,
	}

	return service, nil
}

// GetEmailChannel returns the email notification channel
func (k *kafkaServiceImpl) GetEmailChannel() chan string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.emailChannel
}

// GetSlackChannel returns the slack notification channel
func (k *kafkaServiceImpl) GetSlackChannel() chan string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.slackChannel
}

// GetIOSPushNotificationChannel returns the iOS push notification channel
func (k *kafkaServiceImpl) GetIOSPushNotificationChannel() chan string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.iosPushNotificationChannel
}

// GetAndroidPushNotificationChannel returns the Android push notification channel
func (k *kafkaServiceImpl) GetAndroidPushNotificationChannel() chan string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.androidPushNotificationChannel
}

// Close closes all channels and marks the service as closed
func (k *kafkaServiceImpl) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.closed {
		return
	}

	k.closed = true

	// Close all channels
	close(k.emailChannel)
	close(k.slackChannel)
	close(k.iosPushNotificationChannel)
	close(k.androidPushNotificationChannel)
}

// getEnvAsInt reads an environment variable and converts it to int
// If the environment variable is not set or cannot be parsed, it returns the default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s: %s. Using default: %d\n", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}
