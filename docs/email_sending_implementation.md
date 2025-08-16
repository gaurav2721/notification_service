# Email Sending Implementation

This document explains how the `EmailNotificationRequest` payload is processed and emails are sent through the notification service.

## Overview

The email sending functionality is implemented through a multi-layered architecture:

1. **HTTP Handler** - Receives notification requests and queues them to Kafka
2. **Kafka Consumer** - Processes queued messages and extracts email data
3. **Email Processor** - Converts notification data to `EmailNotificationRequest` and sends emails
4. **Email Service** - Actually sends emails using SMTP or other email providers

## Architecture Flow

```
HTTP Request → Handler → Kafka Channel → Consumer → Email Processor → Email Service → SMTP
```

## Components

### 1. EmailNotificationRequest Model

The `EmailNotificationRequest` is the core data structure for email notifications:

```go
type EmailNotificationRequest struct {
    ID         string                 `json:"id"`
    Type       string                 `json:"type"`
    Content    map[string]interface{} `json:"content"`
    Recipients []string               `json:"recipients"`
    From       *EmailSender           `json:"from,omitempty"`
}

type EmailSender struct {
    Email string `json:"email"`
}
```

**Required Content Fields:**
- `subject` - Email subject line
- `email_body` - HTML email body content

### 2. HTTP Handler (notification_handlers.go)

The handler receives notification requests and processes them:

```go
// SendNotification handles POST /notifications
func (h *NotificationHandler) SendNotification(c *gin.Context) {
    // 1. Parse request
    // 2. Validate "from" field for email notifications
    // 3. Get recipient information from user service
    // 4. Create personalized notification messages
    // 5. Post to Kafka channels
}
```

**Key Features:**
- Validates that email notifications include a "from" field
- Fetches user information to personalize content
- Creates JSON payloads for Kafka channels
- Handles multiple recipients

### 3. Kafka Integration

Messages are sent to Kafka channels based on notification type:

```go
func (h *NotificationHandler) postToKafkaChannel(notificationType string, message map[string]interface{}) error {
    switch notificationType {
    case "email":
        h.kafkaService.GetEmailChannel() <- messageStr
    // ... other notification types
    }
}
```

### 4. Email Processor (email_processor.go)

The email processor is responsible for:

1. **Parsing Kafka messages** - Extracts notification data from JSON payload
2. **Data validation** - Ensures required fields are present
3. **EmailNotificationRequest creation** - Converts parsed data to the proper format
4. **Email sending** - Uses the email service to send actual emails

```go
func (ep *emailProcessor) ProcessNotification(ctx context.Context, message NotificationMessage) error {
    // 1. Parse JSON payload
    var notificationData map[string]interface{}
    json.Unmarshal([]byte(message.Payload), &notificationData)
    
    // 2. Extract recipient email
    recipientData := notificationData["recipient"].(map[string]interface{})
    email := recipientData["email"].(string)
    
    // 3. Extract content
    content := notificationData["content"].(map[string]interface{})
    
    // 4. Extract "from" information
    fromEmail := extractFromEmail(notificationData)
    
    // 5. Create EmailNotificationRequest
    emailNotification := &models.EmailNotificationRequest{
        ID:         message.ID,
        Type:       string(message.Type),
        Content:    content,
        Recipients: []string{email},
        From:       &models.EmailSender{Email: fromEmail},
    }
    
    // 6. Send email
    response, err := ep.emailService.SendEmail(ctx, emailNotification)
    
    return err
}
```

### 5. Email Service (email_service.go)

The email service handles actual email sending:

```go
func (es *EmailServiceImpl) SendEmail(ctx context.Context, notification interface{}) (interface{}, error) {
    // 1. Type assertion and validation
    notif := notification.(*models.EmailNotificationRequest)
    models.ValidateEmailNotification(notif)
    
    // 2. Create email message
    m := gomail.NewMessage()
    m.SetHeader("From", fromEmail)
    m.SetHeader("To", notif.Recipients...)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)
    
    // 3. Send email
    es.dialer.DialAndSend(m)
    
    // 4. Return response
    return &models.EmailResponse{...}, nil
}
```

## Configuration

### Environment Variables

The email service uses these environment variables:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### Fallback to Mock Service

If SMTP configuration is missing, the service automatically falls back to a mock service for testing:

```go
func NewEmailService() EmailService {
    host := os.Getenv("SMTP_HOST")
    // ... check other env vars
    
    if host == "" || portStr == "" || username == "" || password == "" {
        return NewMockEmailService() // Fallback to mock
    }
    
    // Create real SMTP service
    return &EmailServiceImpl{dialer: dialer}
}
```

## Usage Examples

### 1. Sending Email via HTTP API

```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Welcome!",
      "email_body": "<h1>Welcome to our service!</h1>"
    },
    "recipients": ["user@example.com"],
    "from": {
      "email": "noreply@example.com"
    }
  }'
```

### 2. Programmatic Usage

```go
// Create email notification
notification := &models.EmailNotificationRequest{
    ID:         "email-123",
    Type:       "email",
    Content: map[string]interface{}{
        "subject":     "Welcome!",
        "email_body":  "<h1>Welcome!</h1>",
    },
    Recipients: []string{"user@example.com"},
    From: &models.EmailSender{
        Email: "noreply@example.com",
    },
}

// Send email
emailService := email.NewEmailService()
response, err := emailService.SendEmail(ctx, notification)
```

### 3. Custom Email Service Integration

```go
// Create custom email service (SendGrid, AWS SES, etc.)
customService := &MyEmailService{}

// Create processor with custom service
processor := consumers.NewEmailProcessorWithService(customService)

// Process notification
message := consumers.NotificationMessage{...}
err := processor.ProcessNotification(ctx, message)
```

## Error Handling

The implementation includes comprehensive error handling:

1. **Validation Errors** - Missing required fields, invalid email addresses
2. **SMTP Errors** - Connection failures, authentication issues
3. **JSON Parsing Errors** - Malformed notification payloads
4. **Recipient Errors** - Missing or invalid recipient information

## Testing

The implementation includes comprehensive tests:

```bash
# Run email processor tests
go test ./external_services/consumers/ -v

# Run email service tests
go test ./external_services/email/ -v

# Run example
go run examples/email_sending_example.go
```

## Monitoring and Logging

The implementation includes detailed logging:

- **Debug logs** - Processing steps and data flow
- **Info logs** - Successful email sending with metadata
- **Error logs** - Failed operations with error details
- **Structured logging** - JSON-formatted logs with fields

## Performance Considerations

1. **Async Processing** - Emails are processed asynchronously via Kafka
2. **Worker Pools** - Multiple workers can process emails concurrently
3. **Connection Pooling** - SMTP connections are reused
4. **Batch Processing** - Multiple emails can be processed in batches

## Security Considerations

1. **Email Validation** - All email addresses are validated
2. **SMTP Authentication** - Secure SMTP authentication
3. **Content Sanitization** - HTML content should be sanitized
4. **Rate Limiting** - Implement rate limiting to prevent abuse

## Future Enhancements

1. **Template Support** - HTML email templates
2. **Attachment Support** - File attachments
3. **Delivery Tracking** - Email delivery status tracking
4. **Retry Logic** - Automatic retry for failed emails
5. **Multiple Providers** - Support for multiple email providers
6. **Analytics** - Email open/click tracking 