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
		EmailWorkerCount:       getEnvAsInt(constants.EmailWorkerCountEnvVar, constants.DefaultEmailWorkerCount),
		SlackWorkerCount:       getEnvAsInt(constants.SlackWorkerCountEnvVar, constants.DefaultSlackWorkerCount),
		IOSPushWorkerCount:     getEnvAsInt(constants.IOSPushWorkerCountEnvVar, constants.DefaultIOSPushWorkerCount),
		AndroidPushWorkerCount: getEnvAsInt(constants.AndroidPushWorkerCountEnvVar, constants.DefaultAndroidPushWorkerCount),
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
