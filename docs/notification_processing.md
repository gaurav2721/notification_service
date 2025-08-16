# Enhanced Notification Processing System

## Overview

The notification service now includes an enhanced processing system that intelligently handles different notification types by fetching relevant user information and formatting messages appropriately for each channel.

## Key Features

### 1. Type-Based User Information Fetching

The system now fetches user information based on the notification type:

#### **Email Notifications**
- **Fetches**: User's email address
- **Recipient**: Email address from `UserNotificationInfo.Email`
- **Channel**: Email Kafka channel
- **Validation**: Ensures user has a valid email address

#### **Slack Notifications**
- **Fetches**: User's Slack channel information
- **Recipient**: Slack channel from `UserNotificationInfo.SlackChannel`
- **Channel**: Slack Kafka channel
- **Validation**: Ensures user has a configured Slack channel

#### **In-App Notifications**
- **Fetches**: User's device information (tokens, OS, device type)
- **Recipient**: Device tokens grouped by platform (iOS/Android)
- **Channel**: iOS push or Android push Kafka channels
- **Logic**: Automatically determines whether to use iOS or Android push based on user's devices

#### **Direct Push Notifications (ios_push/android_push)**
- **Fetches**: User's device information filtered by platform
- **Recipient**: Device tokens for the specific platform
- **Channel**: Platform-specific Kafka channel
- **Validation**: Ensures user has active devices for the specified platform

### 2. Intelligent Channel Routing

The system automatically routes notifications to the appropriate channels:

```go
// For in_app notifications
if len(iosDevices) > 0 {
    // Route to iOS push channel
    h.postToKafkaChannel("ios_push", iosMessage)
}

if len(androidDevices) > 0 {
    // Route to Android push channel
    h.postToKafkaChannel("android_push", androidMessage)
}
```

### 3. Proper Message Formatting

Each notification type gets properly formatted messages for consumers:

#### **Email Message Format**
```json
{
  "notification_id": "1234567890",
  "type": "email",
  "content": {
    "subject": "Welcome to Our Service",
    "email_body": "Hi John Doe, welcome to our platform!",
    "recipient_name": "John Doe",
    "recipient_email": "john.doe@example.com"
  },
  "template": null,
  "created_at": "2024-01-15T10:00:00Z",
  "recipient": {
    "user_id": "user-123",
    "email": "john.doe@example.com",
    "full_name": "John Doe"
  },
  "from": {
    "email": "noreply@company.com"
  }
}
```

#### **Slack Message Format**
```json
{
  "notification_id": "1234567890",
  "type": "slack",
  "content": {
    "text": "Hi John Doe, you have a new message!",
    "recipient_name": "John Doe",
    "recipient_slack_id": "U1234567890"
  },
  "template": null,
  "created_at": "2024-01-15T10:00:00Z",
  "recipient": {
    "user_id": "user-123",
    "slack_user_id": "U1234567890",
    "slack_channel": "#general",
    "full_name": "John Doe"
  }
}
```

#### **Push Notification Message Format**
```json
{
  "notification_id": "1234567890",
  "type": "ios_push",
  "content": {
    "title": "New Message",
    "body": "Hi John Doe, you have a new message!",
    "recipient_name": "John Doe"
  },
  "template": null,
  "created_at": "2024-01-15T10:00:00Z",
  "recipient": {
    "user_id": "user-123",
    "full_name": "John Doe",
    "device_tokens": ["device_token_1", "device_token_2"],
    "device_count": 2
  },
  "devices": [
    {
      "id": "device-1",
      "device_token": "device_token_1",
      "device_type": "ios",
      "app_version": "1.0.0",
      "os_version": "iOS 16.0",
      "device_model": "iPhone 14"
    }
  ]
}
```

## Implementation Details

### 1. User Information Fetching

The system uses `GetUserNotificationInfo` to fetch comprehensive user data:

```go
userNotificationInfo, err := h.userService.GetUserNotificationInfo(c.Request.Context(), user.ID)
if err != nil {
    logrus.WithError(err).Error("Failed to get user notification info")
    continue
}
```

### 2. Device Filtering and Grouping

For push notifications, devices are filtered and grouped by platform:

```go
// Group devices by type
iosDevices := make([]*models.UserDeviceInfo, 0)
androidDevices := make([]*models.UserDeviceInfo, 0)

for _, device := range userInfo.Devices {
    if device.IsActive && device.DeviceToken != "" {
        switch device.DeviceType {
        case "ios":
            iosDevices = append(iosDevices, device)
        case "android":
            androidDevices = append(androidDevices, device)
        }
    }
}
```

### 3. Message Creation Methods

Three specialized methods create properly formatted messages:

- `createEmailMessage()` - Email-specific formatting
- `createSlackMessage()` - Slack-specific formatting  
- `createPushMessage()` - Push notification formatting

### 4. Channel Routing

Messages are routed to appropriate Kafka channels:

```go
switch notificationType {
case "email":
    h.kafkaService.GetEmailChannel() <- messageStr
case "slack":
    h.kafkaService.GetSlackChannel() <- messageStr
case "ios_push":
    h.kafkaService.GetIOSPushNotificationChannel() <- messageStr
case "android_push":
    h.kafkaService.GetAndroidPushNotificationChannel() <- messageStr
}
```

## Error Handling

### 1. Missing User Information

The system gracefully handles missing user information:

```go
// Email without email address
if userInfo.Email == "" {
    logrus.WithField("user_id", userInfo.ID).Warn("User has no email address")
    return responses, nil
}

// Slack without channel
if userInfo.SlackChannel == "" {
    logrus.WithField("user_id", userInfo.ID).Warn("User has no slack channel")
    return responses, nil
}

// Push without devices
if len(userInfo.Devices) == 0 {
    logrus.WithField("user_id", userInfo.ID).Warn("User has no active devices")
    return responses, nil
}
```

### 2. Channel Failures

Individual channel failures don't stop processing:

```go
err := h.postToKafkaChannel("ios_push", iosMessage)
if err != nil {
    logrus.WithError(err).Error("Failed to post iOS push notification")
} else {
    // Create success response
}
```

## Response Format

The system returns detailed responses for each successful notification:

```json
{
  "notification_id": "1234567890",
  "total_recipients": 3,
  "queued_count": 5,
  "responses": [
    {
      "id": "1234567890",
      "status": "queued",
      "message": "Email notification queued for user John Doe",
      "sent_at": "2024-01-15T10:00:00Z",
      "channel": "email"
    },
    {
      "id": "1234567890", 
      "status": "queued",
      "message": "iOS push notification queued for user John Doe (2 devices)",
      "sent_at": "2024-01-15T10:00:00Z",
      "channel": "ios_push"
    }
  ]
}
```

## Benefits

### 1. **Intelligent Routing**
- Automatically determines the correct channels based on user data
- Handles multi-device users efficiently
- Routes in_app notifications to appropriate push channels

### 2. **Proper Message Formatting**
- Channel-specific message formats
- Consumer-ready message structures
- Includes all necessary metadata

### 3. **Robust Error Handling**
- Graceful handling of missing user data
- Continues processing other users if one fails
- Detailed logging for debugging

### 4. **Scalable Architecture**
- Type-specific processing logic
- Easy to extend for new notification types
- Clean separation of concerns

## Usage Examples

### Send Email Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Welcome",
      "email_body": "Welcome to our platform!"
    },
    "recipients": ["user-123"],
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

### Send In-App Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "content": {
      "title": "New Message",
      "body": "You have a new message!"
    },
    "recipients": ["user-123"]
  }'
```

The system will automatically:
1. Fetch user's device information
2. Group devices by platform (iOS/Android)
3. Send to appropriate push channels
4. Return responses for each platform 