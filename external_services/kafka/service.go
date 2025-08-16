package kafka

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
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
	logrus.Debug("Creating new Kafka service")

	// Read buffer sizes from environment variables with defaults
	emailBufferSize := getEnvAsInt("EMAIL_CHANNEL_BUFFER_SIZE", 100)
	slackBufferSize := getEnvAsInt("SLACK_CHANNEL_BUFFER_SIZE", 100)
	iosPushBufferSize := getEnvAsInt("IOS_PUSH_CHANNEL_BUFFER_SIZE", 100)
	androidPushBufferSize := getEnvAsInt("ANDROID_PUSH_CHANNEL_BUFFER_SIZE", 100)

	logrus.WithFields(logrus.Fields{
		"email_buffer_size":   emailBufferSize,
		"slack_buffer_size":   slackBufferSize,
		"ios_buffer_size":     iosPushBufferSize,
		"android_buffer_size": androidPushBufferSize,
	}).Debug("Kafka service buffer sizes configured")

	service := &kafkaServiceImpl{
		emailChannel:                   make(chan string, emailBufferSize),
		slackChannel:                   make(chan string, slackBufferSize),
		iosPushNotificationChannel:     make(chan string, iosPushBufferSize),
		androidPushNotificationChannel: make(chan string, androidPushBufferSize),
		closed:                         false,
	}

	logrus.Debug("Kafka service created successfully")
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
		logrus.Debug("Kafka service is already closed")
		return
	}

	logrus.Debug("Closing Kafka service")
	k.closed = true

	// Close all channels
	close(k.emailChannel)
	close(k.slackChannel)
	close(k.iosPushNotificationChannel)
	close(k.androidPushNotificationChannel)

	logrus.Debug("Kafka service closed successfully")
}

// getEnvAsInt reads an environment variable and converts it to int
// If the environment variable is not set or cannot be parsed, it returns the default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		logrus.WithFields(logrus.Fields{
			"key":           key,
			"default_value": defaultValue,
		}).Debug("Environment variable not set, using default")
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":           key,
			"value":         valueStr,
			"default_value": defaultValue,
			"error":         err.Error(),
		}).Warn("Invalid environment variable value, using default")
		return defaultValue
	}

	logrus.WithFields(logrus.Fields{
		"key":   key,
		"value": value,
	}).Debug("Environment variable parsed successfully")
	return value
}
