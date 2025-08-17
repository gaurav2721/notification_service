package constants

// Environment variable constants
const (
	// Server configuration
	PORT = "PORT"

	// Logging
	LOG_LEVEL = "LOG_LEVEL"

	// API Security
	API_KEY = "API_KEY"

	// Feature flags
	ENABLE_USER_ROUTES = "ENABLE_USER_ROUTES"

	// SMTP Configuration
	SMTP_HOST     = "SMTP_HOST"
	SMTP_PORT     = "SMTP_PORT"
	SMTP_USERNAME = "SMTP_USERNAME"
	SMTP_PASSWORD = "SMTP_PASSWORD"

	// FCM Configuration
	FCM_SERVER_KEY = "FCM_SERVER_KEY"
	FCM_TIMEOUT    = "FCM_TIMEOUT"
	FCM_BATCH_SIZE = "FCM_BATCH_SIZE"

	// Slack Configuration
	SLACK_BOT_TOKEN  = "SLACK_BOT_TOKEN"
	SLACK_CHANNEL_ID = "SLACK_CHANNEL_ID"

	// APNS Configuration
	APNS_BUNDLE_ID        = "APNS_BUNDLE_ID"
	APNS_KEY_ID           = "APNS_KEY_ID"
	APNS_TEAM_ID          = "APNS_TEAM_ID"
	APNS_PRIVATE_KEY_PATH = "APNS_PRIVATE_KEY_PATH"
	APNS_TIMEOUT          = "APNS_TIMEOUT"

	// Worker Configuration
	EmailWorkerCountEnvVar       = "EMAIL_WORKER_COUNT"
	SlackWorkerCountEnvVar       = "SLACK_WORKER_COUNT"
	IOSPushWorkerCountEnvVar     = "IOS_PUSH_WORKER_COUNT"
	AndroidPushWorkerCountEnvVar = "ANDROID_PUSH_WORKER_COUNT"

	// Kafka Buffer Configuration
	EmailChannelBufferSizeEnvVar       = "EMAIL_CHANNEL_BUFFER_SIZE"
	SlackChannelBufferSizeEnvVar       = "SLACK_CHANNEL_BUFFER_SIZE"
	IOSPushChannelBufferSizeEnvVar     = "IOS_PUSH_CHANNEL_BUFFER_SIZE"
	AndroidPushChannelBufferSizeEnvVar = "ANDROID_PUSH_CHANNEL_BUFFER_SIZE"
)

// Default values for environment variables
const (
	// Server configuration defaults
	DefaultPort = "8080"

	// SMTP Configuration defaults
	DefaultSMTPPort = 587

	// FCM Configuration defaults
	DefaultFCMTimeout   = 30
	DefaultFCMBatchSize = 100

	// APNS Configuration defaults
	DefaultAPNSTimeout = 30

	// Worker Configuration defaults
	DefaultEmailWorkerCount       = 5
	DefaultSlackWorkerCount       = 3
	DefaultIOSPushWorkerCount     = 3
	DefaultAndroidPushWorkerCount = 3

	// Kafka Buffer Configuration defaults
	DefaultEmailChannelBufferSize       = 100
	DefaultSlackChannelBufferSize       = 100
	DefaultIOSPushChannelBufferSize     = 100
	DefaultAndroidPushChannelBufferSize = 100
)
