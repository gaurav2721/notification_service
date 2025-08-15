package consumers

import (
	"os"
	"strconv"
)

// NewConsumerManagerFromEnv creates a new consumer manager with configuration from environment variables
func NewConsumerManagerFromEnv(kafkaService interface {
	GetEmailChannel() chan string
	GetSlackChannel() chan string
	GetIOSPushNotificationChannel() chan string
	GetAndroidPushNotificationChannel() chan string
	Close()
}) ConsumerManager {
	config := ConsumerConfig{
		EmailWorkerCount:       getEnvAsInt("EMAIL_WORKER_COUNT", 5),
		SlackWorkerCount:       getEnvAsInt("SLACK_WORKER_COUNT", 3),
		IOSPushWorkerCount:     getEnvAsInt("IOS_PUSH_WORKER_COUNT", 3),
		AndroidPushWorkerCount: getEnvAsInt("ANDROID_PUSH_WORKER_COUNT", 3),
		KafkaService:           kafkaService,
	}

	return NewConsumerManager(config)
}

// getEnvAsInt gets an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
