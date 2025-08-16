# Template Processing Implementation

## Overview

This document describes the implementation of template processing functionality in the notification service. The system now supports dynamic content generation from templates with variable substitution.

## Implementation Details

### 1. Core Function: `processTemplateToContent`

**Location**: `handlers/notification_handlers.go`

**Function Signature**:
```go
func (h *NotificationHandler) processTemplateToContent(
    template *models.TemplateData, 
    notificationType string
) (map[string]interface{}, error)
```

**Purpose**: Takes a template and notification type from `models.NotificationRequest` and outputs processed content.

### 2. Function Features

#### Template Retrieval and Validation
- Retrieves template by ID from the template manager
- Validates that template type matches notification type
- Validates that all required variables are provided in the template data

#### Content Generation by Type
- **Email**: Generates `subject` and `email_body` fields
- **Slack**: Generates `text` field
- **In-App**: Generates `title` and `body` fields

#### Variable Substitution
- Supports `{{variable_name}}` format for template variables
- Replaces all variables with actual values from template data
- Handles various data types (strings, numbers, etc.)

### 3. Integration with Notification Handler

The function is integrated into the `SendNotification` handler:

```go
// Process template if provided and generate content
if request.Template != nil {
    logrus.Debug("Processing template to generate content")
    generatedContent, err := h.processTemplateToContent(request.Template, request.Type)
    if err != nil {
        logrus.WithError(err).Error("Failed to process template")
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Template processing failed: %v", err)})
        return
    }

    // Replace or merge the content with generated content
    if request.Content == nil {
        request.Content = generatedContent
    } else {
        // Merge generated content with existing content, giving priority to generated content
        for key, value := range generatedContent {
            request.Content[key] = value
        }
    }

    logrus.WithField("generated_content", generatedContent).Debug("Template content generated and merged")
}
```

### 4. Updated Service Interface

The notification service interface has been updated to include template retrieval:

```go
type NotificationHandler struct {
    notificationService interface {
        SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
        ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
        GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
        CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
        GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
        GetPredefinedTemplates() []*models.Template
        GetTemplateByID(templateID string) (*models.Template, error)  // NEW
    }
    // ... other fields
}
```

## Usage Examples

### 1. Email Template Processing

**Request**:
```json
{
  "type": "email",
  "template": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "data": {
      "name": "John Doe",
      "platform": "Tuskira",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "account_type": "Premium",
      "activation_link": "https://tuskira.com/activate?token=abc123"
    }
  },
  "recipients": ["user-123"],
  "from": {
    "email": "noreply@tuskira.com"
  }
}
```

**Generated Content**:
```json
{
  "subject": "Welcome to Tuskira, John Doe!",
  "email_body": "Hello John Doe,\n\nWelcome to Tuskira! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: johndoe\n- Email: john.doe@example.com\n- Account Type: Premium\n\nPlease click the following link to activate your account:\nhttps://tuskira.com/activate?token=abc123\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe Tuskira Team"
}
```

### 2. Slack Template Processing

**Request**:
```json
{
  "type": "slack",
  "template": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "data": {
      "alert_type": "Database Connection",
      "system_name": "User Service",
      "severity": "Critical",
      "environment": "Production",
      "message": "Database connection timeout after 30 seconds",
      "timestamp": "2024-01-01T10:00:00Z",
      "action_required": "Check database connectivity and restart service if needed",
      "affected_services": "User authentication, profile management",
      "dashboard_link": "https://grafana.company.com/d/user-service"
    }
  },
  "recipients": ["user-123"]
}
```

**Generated Content**:
```json
{
  "text": "🚨 *Database Connection Alert*\n\n*System:* User Service\n*Severity:* Critical\n*Environment:* Production\n*Message:* Database connection timeout after 30 seconds\n*Timestamp:* 2024-01-01T10:00:00Z\n*Action Required:* Check database connectivity and restart service if needed\n\n*Affected Services:* User authentication, profile management\n*Dashboard:* https://grafana.company.com/d/user-service\n\nPlease take immediate action if this is a critical alert."
}
```

### 3. In-App Template Processing

**Request**:
```json
{
  "type": "in_app",
  "template": {
    "id": "550e8400-e29b-41d4-a716-446655440005",
    "data": {
      "order_id": "ORD-2024-001",
      "status": "shipped",
      "item_count": "3",
      "total_amount": "299.99",
      "status_message": "Your order has been shipped and is on its way!",
      "action_button": "Track Order"
    }
  },
  "recipients": ["user-123"]
}
```

**Generated Content**:
```json
{
  "title": "Order #ORD-2024-001 - shipped",
  "body": "Your order has been shipped.\n\n*Order Details:*\n- Items: 3 items\n- Total: $299.99\n- Status: shipped\n\nYour order has been shipped and is on its way!\n\nTrack Order"
}
```

## Error Handling

The function provides comprehensive error handling:

### 1. Template Not Found
```go
return nil, fmt.Errorf("failed to get template: %v", err)
```

### 2. Template Type Mismatch
```go
return nil, fmt.Errorf("template type %s does not match notification type %s", templateObj.Type, notificationType)
```

### 3. Missing Required Variables
```go
return nil, fmt.Errorf("template validation failed: %v", err)
```

### 4. Unsupported Notification Type
```go
return nil, fmt.Errorf("unsupported notification type: %s", notificationType)
```

## Testing Results

The implementation has been thoroughly tested with the following results:

✅ **Email Template Processing**: Successfully generates subject and email body
✅ **Slack Template Processing**: Successfully generates text content
✅ **In-App Template Processing**: Successfully generates title and body
✅ **Variable Substitution**: Correctly replaces `{{variable_name}}` with actual values
✅ **Type Validation**: Ensures template type matches notification type
✅ **Required Variable Validation**: Validates all required variables are provided
✅ **Error Handling**: Properly handles various error scenarios

## Benefits

1. **Dynamic Content**: Generate personalized content based on template variables
2. **Type Safety**: Ensures template type matches notification type
3. **Validation**: Validates required variables before processing
4. **Flexibility**: Supports merging with existing content
5. **Maintainability**: Clean separation of concerns with dedicated function
6. **Extensibility**: Easy to add support for new notification types

## Future Enhancements

1. **Advanced Variable Processing**: Support for conditional logic and loops
2. **Template Caching**: Cache frequently used templates for better performance
3. **Template Versioning**: Support for multiple template versions
4. **Content Validation**: Validate generated content before sending
5. **Template Analytics**: Track template usage and performance metrics

## Conclusion

The template processing implementation provides a robust, type-safe, and flexible solution for generating dynamic notification content. It seamlessly integrates with the existing notification flow while maintaining backward compatibility and providing comprehensive error handling. 