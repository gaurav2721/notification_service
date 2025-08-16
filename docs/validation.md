# Notification Service Validation System

## Overview

The notification service now includes a comprehensive input validation system that ensures all API requests are properly validated before processing. This system provides detailed error messages and prevents invalid data from being processed.

## Validation Features

### 1. Request Structure Validation
- Validates JSON structure and required fields
- Ensures proper data types for all fields
- Checks for missing or malformed data

### 2. Notification Type Validation
Supported notification types:
- `email` - Email notifications
- `slack` - Slack messages
- `ios_push` - iOS push notifications
- `android_push` - Android push notifications
- `in_app` - In-app notifications

### 3. Content Validation
Each notification type has specific content requirements:

#### Email Notifications
- **Required fields:**
  - `content.subject` (max 255 characters)
  - `content.email_body` (max 10,000 characters)
- **Required sender:**
  - `from.email` (valid email format, max 254 characters)

#### Slack Notifications
- **Required fields:**
  - `content.text` (max 3,000 characters)

#### Push Notifications (iOS, Android, In-App)
- **Required fields:**
  - `content.title` (max 255 characters)
  - `content.body` (max 4,000 characters)

### 4. Recipients Validation
- At least one recipient required
- Maximum 1,000 recipients per notification
- Each recipient must be 1-255 characters
- Only alphanumeric characters, hyphens, and underscores allowed
- No empty or whitespace-only recipients

### 5. Template Validation
- Template ID must be a valid UUID format
- Template data must be provided and non-empty
- Either content OR template must be provided, not both

### 6. Scheduling Validation
- Scheduled time cannot be in the past
- Scheduled time cannot be more than 1 year in the future
- Must be a valid ISO 8601 timestamp

### 7. From Field Validation
- Required for email notifications
- Must be a valid email format
- Not allowed for non-email notifications

## Usage Examples

### Valid Email Notification
```json
{
  "type": "email",
  "content": {
    "subject": "Welcome to Our Service",
    "email_body": "Thank you for joining our platform!"
  },
  "recipients": ["user-123", "user-456"],
  "from": {
    "email": "noreply@company.com"
  }
}
```

### Valid Slack Notification
```json
{
  "type": "slack",
  "content": {
    "text": "New message received from John Doe"
  },
  "recipients": ["user-123"]
}
```

### Valid Push Notification
```json
{
  "type": "ios_push",
  "content": {
    "title": "New Message",
    "body": "You have received a new message"
  },
  "recipients": ["user-123"]
}
```

### Valid Template Notification
```json
{
  "type": "email",
  "template": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "data": {
      "name": "John Doe",
      "company": "Acme Corp"
    }
  },
  "recipients": ["user-123"],
  "from": {
    "email": "noreply@company.com"
  }
}
```

### Valid Scheduled Notification
```json
{
  "type": "email",
  "content": {
    "subject": "Reminder",
    "email_body": "Don't forget about the meeting tomorrow"
  },
  "recipients": ["user-123"],
  "scheduled_at": "2024-01-15T10:00:00Z",
  "from": {
    "email": "noreply@company.com"
  }
}
```

## Error Response Format

When validation fails, the API returns a structured error response:

```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "content.subject",
      "message": "email subject is required"
    },
    {
      "field": "from.email",
      "message": "from email is required for email notifications"
    }
  ]
}
```

## Common Validation Errors

### Missing Required Fields
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "type",
      "message": "notification type is required"
    }
  ]
}
```

### Invalid Email Format
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "from.email",
      "message": "invalid email format"
    }
  ]
}
```

### Invalid Recipients
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "recipients[0]",
      "message": "recipient can only contain alphanumeric characters, hyphens, and underscores"
    }
  ]
}
```

### Content and Template Conflict
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "content/template",
      "message": "content and template cannot be provided simultaneously"
    }
  ]
}
```

### Invalid Notification Type
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "type",
      "message": "invalid notification type: sms. Valid types are: email, slack, ios_push, android_push, in_app"
    }
  ]
}
```

### Past Scheduled Time
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "scheduled_at",
      "message": "scheduled time cannot be in the past"
    }
  ]
}
```

## Implementation Details

### Validation Package Structure
```
validation/
├── notification_validator.go    # Main validation logic
├── notification_validator_test.go # Unit tests
└── middleware/
    └── validation_middleware.go # Gin middleware
```

### Key Components

1. **NotificationValidator**: Main validation engine
2. **ValidationError**: Individual error structure
3. **ValidationResult**: Complete validation result
4. **ValidationMiddleware**: Gin middleware for automatic validation

### Integration Points

The validation system is integrated into:
- `handlers/notification_handlers.go` - Direct validation in handlers
- `middleware/validation_middleware.go` - Reusable middleware
- All notification endpoints

## Testing

The validation system includes comprehensive unit tests covering:
- Valid request scenarios
- Invalid request scenarios
- Edge cases
- Boundary conditions

Run tests with:
```bash
go test ./validation/...
```

## Best Practices

1. **Always validate input** before processing
2. **Use structured error responses** for better client experience
3. **Log validation failures** for debugging
4. **Keep validation rules consistent** across endpoints
5. **Update validation rules** when adding new notification types

## Future Enhancements

1. **Custom validation rules** for specific business logic
2. **Rate limiting validation** for recipients
3. **Content sanitization** for security
4. **Internationalization** of error messages
5. **Validation caching** for performance 