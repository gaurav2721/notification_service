# Notification API Implementation

## Overview

The notification API has been implemented to handle the complete flow of sending notifications to multiple recipients. When the API endpoint `http://localhost:8080/api/v1/notifications` is hit, the system performs the following steps:

1. **Create notification messages for each recipient**
2. **Get recipient information from userService**
3. **Post notifications to relevant Kafka channels**

## Architecture

### Components

1. **NotificationHandler** - Handles HTTP requests and orchestrates the notification flow
2. **UserService** - Provides user information and notification preferences
3. **KafkaService** - Manages Kafka channels for different notification types
4. **NotificationManager** - Core notification processing logic

### Dependencies

The NotificationManager now only requires the following dependencies:
- **userService** - Used to get required information by using userIds as reference
- **kafkaService** - Used to get channels on which we have to push notification
- **scheduler** - For scheduling notifications
- **templateManager** - For template management

## API Endpoint

### POST /api/v1/notifications

**Request Body:**
```json
{
  "type": "email|slack|ios_push|android_push|in_app",
  "content": {
    "subject": "Email subject",
    "email_body": "Email content",
    "text": "Slack message",
    "title": "Push notification title",
    "body": "Push notification body"
  },
  "template": {
    "id": "template_id",
    "data": {
      "variable1": "value1"
    }
  },
  "recipients": ["user-001", "user-002"],
  "scheduled_at": "2024-12-31T23:59:59Z"
}
```

**Response:**
```json
{
  "notification_id": "123456789",
  "total_recipients": 2,
  "queued_count": 2,
  "responses": [
    {
      "id": "123456789",
      "status": "queued",
      "message": "Notification queued for user John Doe",
      "sent_at": "2024-01-01T12:00:00Z",
      "channel": "email"
    }
  ]
}
```

## Implementation Flow

### 1. Request Processing

When a POST request is received:

```go
func (h *NotificationHandler) SendNotification(c *gin.Context) {
    // Parse request body
    var request struct {
        Type        string                 `json:"type" binding:"required"`
        Content     map[string]interface{} `json:"content"`
        Template    *models.TemplateData   `json:"template,omitempty"`
        Recipients  []string               `json:"recipients" binding:"required"`
        ScheduledAt *time.Time             `json:"scheduled_at"`
    }
}
```

### 2. Scheduled Notification Handling

If `scheduled_at` is provided, the notification is scheduled for future delivery:

```go
if request.ScheduledAt != nil {
    // Create notification object for scheduling
    notification := &struct{...}{
        ID:          generateID(),
        Type:        request.Type,
        Content:     request.Content,
        Template:    request.Template,
        Recipients:  request.Recipients,
        ScheduledAt: request.ScheduledAt,
    }
    
    // Schedule notification
    response, err := h.notificationService.ScheduleNotification(c.Request.Context(), notification)
}
```

### 3. Recipient Information Retrieval

For immediate notifications, recipient information is retrieved from the userService:

```go
// Get recipient information from userService
users, err := h.userService.GetUsersByIDs(c.Request.Context(), request.Recipients)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get recipient information: %v", err)})
    return
}
```

### 4. Notification Message Creation

For each recipient, a personalized notification message is created:

```go
func (h *NotificationHandler) createNotificationMessage(notificationID string, request struct{...}, user *models.User) map[string]interface{} {
    message := map[string]interface{}{
        "notification_id": notificationID,
        "type":            request.Type,
        "content":         request.Content,
        "template":        request.Template,
        "created_at":      time.Now(),
        "recipient": map[string]interface{}{
            "user_id":    user.ID,
            "email":      user.Email,
            "full_name":  user.FullName,
            "slack_user_id": user.SlackUserID,
            "slack_channel": user.SlackChannel,
            "phone_number":  user.PhoneNumber,
        },
    }
    
    // Add personalized content
    if request.Content != nil {
        personalizedContent := make(map[string]interface{})
        for key, value := range request.Content {
            personalizedContent[key] = value
        }
        
        // Add user-specific personalization
        if user.FullName != "" {
            personalizedContent["recipient_name"] = user.FullName
        }
        if user.Email != "" {
            personalizedContent["recipient_email"] = user.Email
        }
        
        message["content"] = personalizedContent
    }
    
    return message
}
```

### 5. Kafka Channel Posting

Each notification message is posted to the appropriate Kafka channel:

```go
func (h *NotificationHandler) postToKafkaChannel(notificationType string, message map[string]interface{}) error {
    // Convert message to JSON
    messageJSON, err := json.Marshal(message)
    messageStr := string(messageJSON)
    
    // Post to appropriate channel based on notification type
    switch notificationType {
    case "email":
        select {
        case h.kafkaService.GetEmailChannel() <- messageStr:
            // Message sent successfully
        default:
            return fmt.Errorf("email channel is full")
        }
    case "slack":
        select {
        case h.kafkaService.GetSlackChannel() <- messageStr:
            // Message sent successfully
        default:
            return fmt.Errorf("slack channel is full")
        }
    case "ios_push":
        select {
        case h.kafkaService.GetIOSPushNotificationChannel() <- messageStr:
            // Message sent successfully
        default:
            return fmt.Errorf("iOS push notification channel is full")
        }
    case "android_push":
        select {
        case h.kafkaService.GetAndroidPushNotificationChannel() <- messageStr:
            // Message sent successfully
        default:
            return fmt.Errorf("Android push notification channel is full")
        }
    case "in_app":
        // For in-app notifications, using iOS channel
        select {
        case h.kafkaService.GetIOSPushNotificationChannel() <- messageStr:
            // Message sent successfully
        default:
            return fmt.Errorf("in-app notification channel is full")
        }
    default:
        return fmt.Errorf("unsupported notification type: %s", notificationType)
    }
    
    return nil
}
```

## Supported Notification Types

1. **Email** - Sent to email channel
2. **Slack** - Sent to Slack channel
3. **iOS Push** - Sent to iOS push notification channel
4. **Android Push** - Sent to Android push notification channel
5. **In-App** - Sent to in-app notification channel (currently using iOS channel)

## Error Handling

- **Invalid recipients**: Returns 400 Bad Request if no valid recipients found
- **User service errors**: Returns 500 Internal Server Error if user information cannot be retrieved
- **Kafka channel full**: Logs error but continues with other recipients
- **Unsupported notification type**: Returns error for unsupported types

## Testing

A test script `test_api.sh` is provided to test the API functionality:

```bash
./test_api.sh
```

The script tests:
1. Email notifications
2. Slack notifications
3. iOS push notifications
4. Scheduled notifications
5. Notification status retrieval

## Configuration

The service uses environment variables for configuration. See `env.example` for available options.

## Dependencies

- **Gin** - HTTP framework
- **Kafka** - Message queuing
- **User Service** - User information management
- **Template Manager** - Notification templates

## Future Enhancements

1. **Template Processing** - Full template variable substitution
2. **Retry Logic** - Retry failed notifications
3. **Metrics** - Notification delivery metrics
4. **Rate Limiting** - Prevent notification spam
5. **Bulk Operations** - Optimize for large recipient lists 