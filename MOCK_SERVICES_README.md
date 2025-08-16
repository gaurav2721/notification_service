# Mock Notification Services

This document describes the mock implementations for all notification services (Email, Slack, APNS, FCM) that write notifications to files instead of sending actual notifications.

## Overview

When environment variables are not properly configured, the services automatically fall back to mock implementations that write notification data to files in the `output/` directory.

## Mock Service Files

Each service writes to its own file:

- **Email**: `output/email.txt`
- **Slack**: `output/slack.txt`
- **APNS**: `output/apns.txt`
- **FCM**: `output/fcm.txt`

## Environment Variables

### Email Service
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USERNAME` - SMTP username
- `SMTP_PASSWORD` - SMTP password

### Slack Service
- `SLACK_BOT_TOKEN` - Slack bot token
- `SLACK_CHANNEL_ID` - Default Slack channel ID

### APNS Service (Apple Push Notification Service)
- `APNS_BUNDLE_ID` - App bundle identifier
- `APNS_KEY_ID` - APNS key ID
- `APNS_TEAM_ID` - Apple team ID
- `APNS_PRIVATE_KEY_PATH` - Path to private key file
- `APNS_TIMEOUT` - Request timeout in seconds (optional, default: 30)

### FCM Service (Firebase Cloud Messaging)
- `FCM_SERVER_KEY` - FCM server key
- `FCM_TIMEOUT` - Request timeout in seconds (optional, default: 30)
- `FCM_BATCH_SIZE` - Batch size for sending notifications (optional, default: 1000)

## How to Test Mock Services

### 1. Set up environment variables (empty for mock mode)
Create a `.env` file with empty values for the notification services:

```bash
# Email Configuration (empty for mock mode)
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=

# Slack Configuration (empty for mock mode)
SLACK_BOT_TOKEN=
SLACK_CHANNEL_ID=

# APNS Configuration (empty for mock mode)
APNS_BUNDLE_ID=
APNS_KEY_ID=
APNS_TEAM_ID=
APNS_PRIVATE_KEY_PATH=
APNS_TIMEOUT=30

# FCM Configuration (empty for mock mode)
FCM_SERVER_KEY=
FCM_TIMEOUT=30
FCM_BATCH_SIZE=1000
```

### 2. Run the application
```bash
go run main.go
```

### 3. Send test notifications
Use the API endpoints to send notifications:

```bash
# Test email notification
curl -X POST http://localhost:8080/api/v1/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test_email_001",
    "type": "welcome",
    "content": {
      "subject": "Welcome!",
      "email_body": "<h1>Welcome to our service!</h1>"
    },
    "recipients": ["user@example.com"]
  }'

# Test Slack notification
curl -X POST http://localhost:8080/api/v1/notifications/slack \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test_slack_001",
    "type": "alert",
    "content": {
      "text": "🚨 System alert: High CPU usage detected"
    },
    "recipients": ["#alerts"]
  }'

# Test APNS notification
curl -X POST http://localhost:8080/api/v1/notifications/push/ios \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test_apns_001",
    "type": "push",
    "content": {
      "title": "New Message",
      "body": "You have a new message"
    },
    "recipients": ["ios_device_token_123"]
  }'

# Test FCM notification
curl -X POST http://localhost:8080/api/v1/notifications/push/android \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test_fcm_001",
    "type": "push",
    "content": {
      "title": "App Update",
      "body": "New version available"
    },
    "recipients": ["android_device_token_123"]
  }'
```

### 4. Check output files
After sending notifications, check the `output/` directory for the generated files:

```bash
ls -la output/
cat output/email.txt
cat output/slack.txt
cat output/apns.txt
cat output/fcm.txt
```

## File Format

Each notification is written to the file in JSON format with a separator:

```
=== EMAIL NOTIFICATION ===
{
  "timestamp": "2024-01-15T10:30:00Z",
  "id": "test_email_001",
  "type": "welcome",
  "content": {
    "subject": "Welcome!",
    "email_body": "<h1>Welcome to our service!</h1>"
  },
  "recipients": ["user@example.com"],
  "template": null,
  "status": "mock_sent",
  "channel": "email"
}

```

## Switching to Real Services

To use real services instead of mocks, populate the environment variables with actual values:

```bash
# For real email service
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# For real Slack service
SLACK_BOT_TOKEN=xoxb-your-slack-bot-token
SLACK_CHANNEL_ID=C1234567890

# For real APNS service
APNS_BUNDLE_ID=com.yourcompany.yourapp
APNS_KEY_ID=your-key-id
APNS_TEAM_ID=your-team-id
APNS_PRIVATE_KEY_PATH=/path/to/your/private-key.p8

# For real FCM service
FCM_SERVER_KEY=your-fcm-server-key
```

## Implementation Details

### Service Initializers

Each service has an initializer function that checks environment variables:

- `email.NewEmailService()` - Checks SMTP_* variables
- `slack.NewSlackService()` - Checks SLACK_* variables  
- `apns.NewAPNSService()` - Checks APNS_* variables
- `fcm.NewFCMService()` - Checks FCM_* variables

If any required environment variable is missing or empty, the service returns a mock implementation.

### Mock Service Structure

Each mock service:
1. Implements the same interface as the real service
2. Writes notification data to a JSON file
3. Returns a mock response indicating success
4. Creates the output directory if it doesn't exist

### Error Handling

Mock services handle errors gracefully:
- File I/O errors are returned as errors
- Invalid notification payloads return appropriate error types
- The output directory is created automatically if missing 