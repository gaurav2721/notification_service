package kafka

// KafkaService represents a generic Kafka service interface
type KafkaService interface {
	GetEmailChannel() chan string
	GetSlackChannel() chan string
	GetIOSPushNotificationChannel() chan string
	GetAndroidPushNotificationChannel() chan string
	Close()
}
