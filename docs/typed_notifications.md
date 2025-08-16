# Typed Notification System

## Overview

The notification service now uses strongly-typed notification request structures instead of generic maps. This provides better type safety, validation, and ensures that messages are properly formatted for each notification channel.

## Typed Notification Request Structures

### 1. Email Notifications

**Request Type**: `models.EmailNotificationRequest`

```go
type EmailNotificationRequest struct {
    ID         string       `json:"id"`
    Type       string       `json:"type"`
    Content    EmailContent `json:"content"`
    Recipients []string     `json:"recipients"`
    From       *EmailSender `json:"from,omitempty"`
}

type EmailContent struct {
    Subject   string `json:"subject"`
    EmailBody string `json:"email_body"`
}

type EmailSender struct {
    Email string `json:"email"`
}
```

**Example Message**:
```json
{
  "id": "1234567890",
  "type": "email",
  "content": {
    "subject": "Welcome to Our Service",
    "email_body": "Hi John Doe, welcome to our platform!"
  },
  "recipients": ["john.doe@example.com"],
  "from": {
    "email": "noreply@company.com"
  }
}
```

### 2. Slack Notifications

**Request Type**: `models.SlackNotificationRequest`

```go
type SlackNotificationRequest struct {
    ID         string       `json:"id"`
    Type       string       `json:"type"`
    Content    SlackContent `json:"content"`
    Recipients []string     `json:"recipients"`
}

type SlackContent struct {
    Text string `json:"text"`
}
```

**Example Message**:
```json
{
  "id": "1234567890",
  "type": "slack",
  "content": {
    "text": "Hi John Doe, you have a new message!"
  },
  "recipients": ["#general"]
}
```

### 3. iOS Push Notifications

**Request Type**: `models.APNSNotificationRequest`

```go
type APNSNotificationRequest struct {
    ID         string      `json:"id"`
    Type       string      `json:"type"`
    Content    APNSContent `json:"content"`
    Recipients []string    `json:"recipients"`
}

type APNSContent struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}
```

**Example Message**:
```json
{
  "id": "1234567890",
  "type": "ios_push",
  "content": {
    "title": "New Message",
    "body": "Hi John Doe, you have a new message!"
  },
  "recipients": ["device_token_1", "device_token_2"]
}
```

### 4. Android Push Notifications

**Request Type**: `models.FCMNotificationRequest`

```go
type FCMNotificationRequest struct {
    ID         string     `json:"id"`
    Type       string     `json:"type"`
    Content    FCMContent `json:"content"`
    Recipients []string   `json:"recipients"`
}

type FCMContent struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}
```

**Example Message**:
```json
{
  "id": "1234567890",
  "type": "android_push",
  "content": {
    "title": "New Message",
    "body": "Hi John Doe, you have a new message!"
  },
  "recipients": ["fcm_device_token_1", "fcm_device_token_2"]
}
```

## Message Creation Methods

### 1. Email Message Creation

```go
func (h *NotificationHandler) createEmailMessage(
    notificationID string, 
    request models.NotificationRequest, 
    userInfo *models.UserNotificationInfo,
) *models.EmailNotificationRequest
```

**Features**:
- Extracts subject and email_body from request content
- Personalizes content with user information using template variables
- Sets recipient to user's email address
- Includes from field if provided in original request

**Template Variables**:
- `{{recipient_name}}` or `{{name}}` → User's full name
- `{{recipient_email}}` or `{{email}}` → User's email address

### 2. Slack Message Creation

```go
func (h *NotificationHandler) createSlackMessage(
    notificationID string, 
    request models.NotificationRequest, 
    userInfo *models.UserNotificationInfo,
) *models.SlackNotificationRequest
```

**Features**:
- Extracts text from request content
- Personalizes content with user information
- Sets recipient to user's Slack channel

**Template Variables**:
- `{{recipient_name}}` or `{{name}}` → User's full name
- `{{recipient_slack_id}}` or `{{slack_id}}` → User's Slack ID

### 3. Push Message Creation

```go
func (h *NotificationHandler) createPushMessage(
    notificationID string, 
    request models.NotificationRequest, 
    userInfo *models.UserNotificationInfo, 
    devices []*models.UserDeviceInfo, 
    pushType string,
) interface{}
```

**Features**:
- Extracts title and body from request content
- Personalizes content with user information
- Returns appropriate typed struct based on pushType:
  - `ios_push` → `*models.APNSNotificationRequest`
  - `android_push` → `*models.FCMNotificationRequest`
- Sets recipients to device tokens

**Template Variables**:
- `{{recipient_name}}` or `{{name}}` → User's full name

## Benefits of Typed Notifications

### 1. **Type Safety**
- Compile-time validation of notification structure
- Prevents runtime errors from malformed messages
- IDE support with autocomplete and error detection

### 2. **Validation**
- Built-in validation methods for each notification type
- Consistent validation rules across the system
- Clear error messages for invalid data

### 3. **Consumer Compatibility**
- Messages are properly formatted for each consumer service
- No need for consumers to parse generic maps
- Direct compatibility with email, Slack, APNS, and FCM services

### 4. **Maintainability**
- Clear structure for each notification type
- Easy to extend with new fields
- Self-documenting code

## Usage Examples

### Creating an Email Notification

```go
// The system automatically creates a typed EmailNotificationRequest
emailMessage := h.createEmailMessage(notificationID, request, userInfo)

// The message is properly typed and validated
if err := models.ValidateEmailNotification(emailMessage); err != nil {
    return err
}

// Send to Kafka channel
err := h.postToKafkaChannel("email", emailMessage)
```

### Creating a Push Notification

```go
// The system automatically determines the correct type
pushMessage := h.createPushMessage(notificationID, request, userInfo, devices, "ios_push")

// Type assertion for specific handling
if apnsMessage, ok := pushMessage.(*models.APNSNotificationRequest); ok {
    if err := models.ValidateAPNSNotification(apnsMessage); err != nil {
        return err
    }
}

// Send to Kafka channel
err := h.postToKafkaChannel("ios_push", pushMessage)
```

## Migration from Generic Maps

### Before (Generic Map)
```go
message := map[string]interface{}{
    "notification_id": notificationID,
    "type":            "email",
    "content":         request.Content,
    "recipient": map[string]interface{}{
        "user_id": userInfo.ID,
        "email":   userInfo.Email,
    },
}
```

### After (Typed Struct)
```go
emailNotification := &models.EmailNotificationRequest{
    ID:   notificationID,
    Type: "email",
    Content: models.EmailContent{
        Subject:   subject,
        EmailBody: emailBody,
    },
    Recipients: []string{userInfo.Email},
}
```

## Validation

Each notification type has built-in validation:

```go
// Email validation
if err := models.ValidateEmailNotification(emailNotification); err != nil {
    return err
}

// Slack validation
if err := models.ValidateSlackNotification(slackNotification); err != nil {
    return err
}

// APNS validation
if err := models.ValidateAPNSNotification(apnsNotification); err != nil {
    return err
}

// FCM validation
if err := models.ValidateFCMNotification(fcmNotification); err != nil {
    return err
}
```

## Error Handling

The system gracefully handles type conversion errors:

```go
// Safe type assertion for push messages
if apnsMessage, ok := pushMessage.(*models.APNSNotificationRequest); ok {
    // Handle APNS message
} else if fcmMessage, ok := pushMessage.(*models.FCMNotificationRequest); ok {
    // Handle FCM message
} else {
    // Handle generic fallback
}
```

## Future Enhancements

1. **Template Support**: Enhanced template variable replacement
2. **Custom Fields**: Support for notification-specific custom fields
3. **Batch Processing**: Optimized handling of multiple notifications
4. **Metrics**: Type-specific metrics and monitoring
5. **Schema Evolution**: Versioned notification schemas 