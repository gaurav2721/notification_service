package consumers

import (
	"os"
	"strconv"

	"github.com/gaurav2721/notification-service/constants"
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
		EmailWorkerCount:       getEnvAsInt(constants.EmailWorkerCountEnvVar, 5),
		SlackWorkerCount:       getEnvAsInt(constants.SlackWorkerCountEnvVar, 3),
		IOSPushWorkerCount:     getEnvAsInt(constants.IOSPushWorkerCountEnvVar, 3),
		AndroidPushWorkerCount: getEnvAsInt(constants.AndroidPushWorkerCountEnvVar, 3),
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
