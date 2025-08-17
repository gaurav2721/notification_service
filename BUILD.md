# Build and Run Instructions

## Docker Commands

### Build and Run with Docker
1. `make docker-build`
2. `make docker-run`

### View Docker Container Output
1. `make docker-exec`

## Local Development

### Build and Run Locally
1. `make build`
2. `make run`

### Debug Mode
1. `make run-debug`

## Testing

### Run Unit Tests
1. `make test`




## Environment Variables (.env)

The notification service uses environment variables for configuration. Create a `.env` file in the root directory with the following variables:

### Server Configuration
```env
# Server port (default: 8080)
PORT=8080

# Logging level (debug, info, warn, error)
LOG_LEVEL=info

# API key for authentication
API_KEY=your-secure-api-key-here

# Enable user routes (false by default)
ENABLE_USER_ROUTES=true
```

### Email Configuration (SMTP)(Optional - If not provided , output will be printed in a text file output/email.txt)
```env
# SMTP server host
SMTP_HOST=smtp.gmail.com

# SMTP server port
SMTP_PORT=587

# SMTP username/email
SMTP_USERNAME=your-email@gmail.com

# SMTP password/app password
SMTP_PASSWORD=your-app-password
```

### Slack Configuration(Optional - If not provided , output will be printed in a text file output/slack.txt)
```env
# Slack bot token
SLACK_BOT_TOKEN=xoxb-your-slack-bot-token

# Slack channel ID
SLACK_CHANNEL_ID=C1234567890
```

### Firebase Cloud Messaging (FCM) Configuration(Optional - If not provided , output will be printed in a text file output/fcm.txt)
```env
# FCM server key
FCM_SERVER_KEY=your-fcm-server-key

# FCM request timeout in seconds
FCM_TIMEOUT=30

# FCM batch size for token processing
FCM_BATCH_SIZE=500
```

### Apple Push Notification Service (APNS) Configuration(Optional - If not provided , output will be printed in a text file output/apns.txt)
```env
# APNS bundle ID
APNS_BUNDLE_ID=com.yourcompany.yourapp

# APNS key ID
APNS_KEY_ID=your-key-id

# APNS team ID
APNS_TEAM_ID=your-team-id

# Path to APNS private key file
APNS_PRIVATE_KEY_PATH=/path/to/AuthKey_XXXXXXXXXX.p8

# APNS request timeout in seconds
APNS_TIMEOUT=30
```

### Kafka Channel Buffer Sizes (Optional)
```env
# Email channel buffer size (default: 100)
EMAIL_CHANNEL_BUFFER_SIZE=100

# Slack channel buffer size (default: 100)
SLACK_CHANNEL_BUFFER_SIZE=100

# iOS push notification channel buffer size (default: 100)
IOS_PUSH_CHANNEL_BUFFER_SIZE=100

# Android push notification channel buffer size (default: 100)
ANDROID_PUSH_CHANNEL_BUFFER_SIZE=100
```

### Default .env File

The notification service comes with a default `.env` file that includes basic configuration. Here's the complete default configuration:

```env
# Server Configuration
PORT=8080
LOG_LEVEL=info
API_KEY=gaurav

# # Email Configuration
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_USERNAME=your-email@gmail.com
# SMTP_PASSWORD=your-app-password

# # Slack Configuration
# SLACK_BOT_TOKEN=xoxb-your-slack-bot-token
# SLACK_CHANNEL_ID=C1234567890

# # Database Configuration (for future use)
# DB_HOST=localhost
# DB_PORT=5432
# DB_NAME=notification_service
# DB_USER=postgres
# DB_PASSWORD=password

# Kafka Channel Buffer Sizes
EMAIL_CHANNEL_BUFFER_SIZE=100
SLACK_CHANNEL_BUFFER_SIZE=100
IOS_PUSH_CHANNEL_BUFFER_SIZE=100
ANDROID_PUSH_CHANNEL_BUFFER_SIZE=100

# Consumer Worker Pool Configuration
EMAIL_WORKER_COUNT=5
SLACK_WORKER_COUNT=3
IOS_PUSH_WORKER_COUNT=3
ANDROID_PUSH_WORKER_COUNT=3 
```

### Setup Instructions
1. The default `.env` file is already included in the repository for ease of testing
2. For production, create a new `.env` file with your actual values

### Security Notes
- The default `.env` file has been committed in this repo for ease of testing
- For production, never commit your `.env` file to version control
- Use strong, unique API keys in production
- Store sensitive credentials securely
- Consider using a secrets management service in production

