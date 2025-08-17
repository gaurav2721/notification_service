# Notification Service API Documentation

## Overview

The Notification Service provides a comprehensive API for sending and managing notifications across multiple channels including email, Slack, and in-app notifications. The service supports both immediate and scheduled delivery with template-based content management.

**Base URL:** `http://localhost:8080`  
**API Version:** `v1`  
**Authentication:** Bearer token required for all endpoints except health checks

## Authentication

All API endpoints require authentication using a Bearer token in the Authorization header:

```
Authorization: Bearer gaurav
```

## API Endpoints

### 1. Send Notification

**Endpoint:** `POST /api/v1/notifications`

Send immediate or scheduled notifications to one or more recipients.

#### Request Body

The request body supports two modes:
1. **Direct Content Mode**: Send notifications with direct content
2. **Template Mode**: Send notifications using predefined or custom templates

##### Direct Content Mode

```json
{
  "type": "email|slack|in_app",
  "content": {
    // Content varies by notification type
  },
  "recipients": ["user-id-1", "user-id-2"],
  "scheduled_at": "2024-01-15T14:00:00Z", // Optional for scheduled notifications
  "from": {
    "email": "noreply@company.com" // Required for email notifications
  }
}
```

##### Template Mode

```json
{
  "type": "email|slack|in_app",
  "template": {
    "id": "template-id",
    "version": 1, // Required - must be a positive integer
    "data": {
      // Template variables
    }
  },
  "recipients": ["user-id-1", "user-id-2"],
  "scheduled_at": "2024-01-15T14:00:00Z", // Optional for scheduled notifications
  "from": {
    "email": "noreply@company.com" // Required for email notifications
  }
}
```

#### Content Structure by Type

##### Email Notifications

```json
{
  "type": "email",
  "content": {
    "subject": "Email Subject",
    "email_body": "Email body content with support for newlines"
  },
  "recipients": ["user-001"],
  "from": {
    "email": "noreply@company.com"
  }
}
```

##### Slack Notifications

```json
{
  "type": "slack",
  "content": {
    "text": "Slack message with *bold* and _italic_ formatting"
  },
  "recipients": ["user-001"]
}
```

##### In-App Notifications

```json
{
  "type": "in_app",
  "content": {
    "title": "Notification Title",
    "body": "Notification body content"
  },
  "recipients": ["user-001"]
}
```

#### Response

**Success Response (200 OK):**
```json
{
  "id": "888e9012-e89b-12d3-a456-426614174020",
  "status": "sent" // or "scheduled" for scheduled notifications
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Invalid request",
  "message": "Detailed error message"
}
```

#### Examples

##### Send Immediate Email
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Welcome to Our Platform!",
      "email_body": "Hello,\n\nWelcome to our platform! We are excited to have you on board.\n\nBest regards,\nThe Team"
    },
    "recipients": ["user-001"],
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

##### Send Scheduled Slack Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "content": {
      "text": "📅 *Daily Standup Reminder*\n\nTime: 9:00 AM\nChannel: #daily-standup\nAgenda: Project updates and blockers"
    },
    "recipients": ["user-001"],
    "scheduled_at": "2024-01-16T09:00:00Z"
  }'
```

### 2. Get Notification Status

**Endpoint:** `GET /api/v1/notifications/{notification_id}`

Retrieve the status of a specific notification by its ID.

#### Path Parameters

- `notification_id` (string, required): The unique identifier of the notification

#### Response

**Success Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "sent" // or "scheduled", "failed", "pending"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Notification not found",
  "message": "Notification with ID 123e4567-e89b-12d3-a456-426614174000 not found"
}
```

#### Example

```bash
curl -X GET http://localhost:8080/api/v1/notifications/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer gaurav"
```

### 3. Get Predefined Templates

**Endpoint:** `GET /api/v1/templates/predefined`

Retrieve all available predefined templates for different notification types.

#### Response

**Success Response (200 OK):**
```json
{
  "count": 7,
  "templates": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Welcome Email Template",
      "type": "email",
      "version": 1,
      "content": {
        "subject": "Welcome to {{platform}}, {{name}}!",
        "email_body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
      },
      "description": "Welcome email template for new user onboarding",
      "required_variables": ["name", "platform", "username", "email", "account_type", "activation_link"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787202379Z"
    }
    // ... more templates
  ]
}
```

#### Example

```bash
curl -X GET http://localhost:8080/api/v1/templates/predefined \
  -H "Authorization: Bearer gaurav"
```

### 4. Create Custom Template

**Endpoint:** `POST /api/v1/templates`

Create a new custom template for notifications.

#### Request Body

```json
{
  "name": "Template Name",
  "type": "email|slack|in_app",
  "content": {
    // Content structure varies by type
  },
  "required_variables": ["var1", "var2"],
  "description": "Template description"
}
```

#### Response

**Success Response (201 Created):**
```json
{
  "id": "template-password-reset-custom",
  "name": "Password Reset Template",
  "type": "email",
  "version": 1,
  "status": "created",
  "created_at": "2025-08-15T18:25:00Z"
}
```

#### Example

```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Password Reset Template",
    "type": "email",
    "content": {
      "subject": "Password Reset Request - {{platform_name}}",
      "email_body": "Hello {{user_name}},\n\nWe received a request to reset your password for your {{platform_name}} account.\n\nTo reset your password, click the link below:\n{{reset_link}}\n\nThis link will expire in {{expiry_hours}} hours.\n\nIf you did not request a password reset, please ignore this email or contact support if you have concerns.\n\nBest regards,\nThe {{platform_name}} Team\n\n---\nThis is an automated message, please do not reply to this email."
    },
    "required_variables": ["user_name", "platform_name", "reset_link", "expiry_hours"],
    "description": "Email template for password reset requests"
  }'
```

### 5. Health Check

**Endpoint:** `GET /health`

Check the health status of the notification service. This endpoint does not require authentication.

#### Response

**Success Response (200 OK):**
```json
{
  "service": "notification-service",
  "status": "healthy",
  "timestamp": "2025-08-15T18:23:52.426265799Z"
}
```

#### Example

```bash
curl -X GET http://localhost:8080/health
```

## Preloaded Info

### User


```json
"users": [
    {
      "id": "user-001",
      "email": "john.doe@company.com",
      "full_name": "John Doe",
      "slack_user_id": "U1234567890",
      "slack_channel": "#general",
      "phone_number": "+1-555-0101",
      "is_active": true,
      "created_at": "2025-02-15T18:23:46.787176921Z",
      "updated_at": "2025-08-15T18:23:46.787197838Z"
    },
    {
      "id": "user-002",
      "email": "jane.smith@company.com",
      "full_name": "Jane Smith",
      "slack_user_id": "U0987654321",
      "slack_channel": "#design",
      "phone_number": "+1-555-0102",
      "is_active": true,
      "created_at": "2025-04-15T18:23:46.787197879Z",
      "updated_at": "2025-08-15T18:23:46.787197963Z"
    },
    {
      "id": "user-003",
      "email": "mike.johnson@company.com",
      "full_name": "Mike Johnson",
      "slack_user_id": "U1122334455",
      "slack_channel": "#marketing",
      "phone_number": "+1-555-0103",
      "is_active": true,
      "created_at": "2024-12-15T18:23:46.787198004Z",
      "updated_at": "2025-08-15T18:23:46.787198088Z"
    },
    {
      "id": "user-004",
      "email": "sarah.wilson@company.com",
      "full_name": "Sarah Wilson",
      "slack_user_id": "U5566778899",
      "slack_channel": "#sales",
      "phone_number": "+1-555-0104",
      "is_active": true,
      "created_at": "2025-06-15T18:23:46.787198129Z",
      "updated_at": "2025-08-15T18:23:46.787198213Z"
    },
    {
      "id": "user-005",
      "email": "david.brown@company.com",
      "full_name": "David Brown",
      "slack_user_id": "U9988776655",
      "slack_channel": "#engineering",
      "phone_number": "+1-555-0105",
      "is_active": true,
      "created_at": "2024-08-15T18:23:46.787198254Z",
      "updated_at": "2025-08-15T18:23:46.787198296Z"
    },
    {
      "id": "user-006",
      "email": "lisa.garcia@company.com",
      "full_name": "Lisa Garcia",
      "slack_user_id": "U4433221100",
      "slack_channel": "#marketing",
      "phone_number": "+1-555-0106",
      "is_active": true,
      "created_at": "2024-10-15T18:23:46.787198338Z",
      "updated_at": "2025-08-15T18:23:46.787198421Z"
    },
    {
      "id": "user-007",
      "email": "robert.taylor@company.com",
      "full_name": "Robert Taylor",
      "slack_user_id": "U1122334455",
      "slack_channel": "#sales",
      "phone_number": "+1-555-0107",
      "is_active": true,
      "created_at": "2024-11-15T18:23:46.787198463Z",
      "updated_at": "2025-08-15T18:23:46.787198546Z"
    },
    {
      "id": "user-008",
      "email": "emma.davis@company.com",
      "full_name": "Emma Davis",
      "slack_user_id": "U6677889900",
      "slack_channel": "#executives",
      "phone_number": "+1-555-0108",
      "is_active": true,
      "created_at": "2024-05-15T18:23:46.787198588Z",
      "updated_at": "2025-08-15T18:23:46.787198629Z"
    }
  ]
```

### UserDeviceInfo

```json
"devices": [
  {
    "id": "device-001",
    "user_id": "user-001",
    "device_token": "ios_token_123456789",
    "device_type": "ios",
    "app_version": "1.2.3",
    "os_version": "iOS 16.0",
    "device_model": "iPhone 14",
    "is_active": true,
    "last_used_at": "2025-08-15T18:23:46.787176921Z",
    "created_at": "2025-06-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T18:23:46.787176921Z"
  },
  {
    "id": "device-002",
    "user_id": "user-001",
    "device_token": "android_token_987654321",
    "device_type": "android",
    "app_version": "1.2.3",
    "os_version": "Android 13",
    "device_model": "Samsung Galaxy S23",
    "is_active": true,
    "last_used_at": "2025-08-15T17:23:46.787176921Z",
    "created_at": "2025-07-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T17:23:46.787176921Z"
  },
  {
    "id": "device-004",
    "user_id": "user-003",
    "device_token": "ios_token_789123456",
    "device_type": "ios",
    "app_version": "1.2.2",
    "os_version": "iOS 15.5",
    "device_model": "iPhone 13",
    "is_active": false,
    "last_used_at": "2025-08-08T18:23:46.787176921Z",
    "created_at": "2025-02-15T18:23:46.787176921Z",
    "updated_at": "2025-08-08T18:23:46.787176921Z"
  },
  {
    "id": "device-005",
    "user_id": "user-004",
    "device_token": "android_token_555666777",
    "device_type": "android",
    "app_version": "1.2.4",
    "os_version": "Android 14",
    "device_model": "Google Pixel 8",
    "is_active": true,
    "last_used_at": "2025-08-15T18:08:46.787176921Z",
    "created_at": "2025-07-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T18:08:46.787176921Z"
  },
  {
    "id": "device-006",
    "user_id": "user-005",
    "device_token": "ios_token_111222333",
    "device_type": "ios",
    "app_version": "1.2.3",
    "os_version": "iOS 17.0",
    "device_model": "iPhone 15 Pro",
    "is_active": true,
    "last_used_at": "2025-08-15T16:23:46.787176921Z",
    "created_at": "2025-05-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T16:23:46.787176921Z"
  },
  {
    "id": "device-008",
    "user_id": "user-007",
    "device_token": "android_token_777888999",
    "device_type": "android",
    "app_version": "1.2.1",
    "os_version": "Android 12",
    "device_model": "OnePlus 9",
    "is_active": true,
    "last_used_at": "2025-08-15T15:23:46.787176921Z",
    "created_at": "2025-04-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T15:23:46.787176921Z"
  },
  {
    "id": "device-009",
    "user_id": "user-008",
    "device_token": "ios_token_000111222",
    "device_type": "ios",
    "app_version": "1.2.5",
    "os_version": "iOS 17.2",
    "device_model": "iPad Pro",
    "is_active": true,
    "last_used_at": "2025-08-15T18:13:46.787176921Z",
    "created_at": "2025-07-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T18:13:46.787176921Z"
  },
  {
    "id": "device-010",
    "user_id": "user-002",
    "device_token": "android_token_333444555",
    "device_type": "android",
    "app_version": "1.2.3",
    "os_version": "Android 13",
    "device_model": "Samsung Galaxy Tab S9",
    "is_active": true,
    "last_used_at": "2025-08-15T14:23:46.787176921Z",
    "created_at": "2025-06-15T18:23:46.787176921Z",
    "updated_at": "2025-08-15T14:23:46.787176921Z"
  }
]
```

### Template

```json
"templates": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Welcome Email Template",
      "type": "email",
      "version": 1,
      "content": {
        "subject": "Welcome to {{platform}}, {{name}}!",
        "email_body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
      },
      "description": "Welcome email template for new user onboarding",
      "required_variables": ["name", "platform", "username", "email", "account_type", "activation_link"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787202379Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Password Reset Template",
      "type": "email",
      "version": 1,
      "content": {
        "subject": "Password Reset Request - {{platform}}",
        "email_body": "Hello {{name}},\n\nWe received a request to reset your password for your {{platform}} account.\n\nTo reset your password, click the link below:\n{{reset_link}}\n\nThis link will expire in {{expiry_hours}} hours.\n\nIf you did not request a password reset, please ignore this email or contact support if you have concerns.\n\nBest regards,\nThe {{platform}} Team"
      },
      "description": "Password reset email template",
      "required_variables": ["name", "platform", "reset_link", "expiry_hours"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787202671Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "Order Confirmation Template",
      "type": "email",
      "version": 1,
      "content": {
        "subject": "Order Confirmation - {{order_id}}",
        "email_body": "Hello {{customer_name}},\n\nThank you for your order! Your order has been confirmed and is being processed.\n\nOrder Details:\n- Order ID: {{order_id}}\n- Order Date: {{order_date}}\n- Total Amount: {{total_amount}}\n- Payment Method: {{payment_method}}\n\nItems:\n{{items_list}}\n\nShipping Address:\n{{shipping_address}}\n\nExpected Delivery: {{delivery_date}}\n\nTrack your order: {{tracking_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
      },
      "description": "Order confirmation email template",
      "required_variables": ["customer_name", "order_id", "order_date", "total_amount", "payment_method", "items_list", "shipping_address", "delivery_date", "tracking_link", "platform"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787202879Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "System Alert Template",
      "type": "slack",
      "version": 1,
      "content": {
        "text": "🚨 *{{alert_type}} Alert*\n\n*System:* {{system_name}}\n*Severity:* {{severity}}\n*Environment:* {{environment}}\n*Message:* {{message}}\n*Timestamp:* {{timestamp}}\n*Action Required:* {{action_required}}\n*Affected Services:* {{affected_services}}\n\n<{{dashboard_link}}|View Dashboard>"
      },
      "description": "Slack alert template for system monitoring",
      "required_variables": ["alert_type", "system_name", "severity", "environment", "message", "timestamp", "action_required", "affected_services", "dashboard_link"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787203338Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "name": "Deployment Notification Template",
      "type": "slack",
      "version": 1,
      "content": {
        "text": "🚀 *Deployment {{status}}*\n\n*Service:* {{service_name}}\n*Environment:* {{environment}}\n*Version:* {{version}}\n*Deployed By:* {{deployed_by}}\n*Duration:* {{duration}}\n\n*Changes Summary:*\n{{changes_summary}}\n\n*Rollback Command:*\n```{{rollback_command}}```\n\n<{{monitoring_link}}|Monitor Service>"
      },
      "description": "Slack notification template for deployment events",
      "required_variables": ["status", "service_name", "environment", "version", "deployed_by", "duration", "changes_summary", "rollback_command", "monitoring_link"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787203504Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "name": "Order Status Update Template",
      "type": "in_app",
      "version": 1,
      "content": {
        "title": "Order #{{order_id}} - {{status}}",
        "body": "Your order with {{item_count}} items ({{total_amount}}) has been {{status}}.\n\n{{status_message}}\n\nTap to {{action_button}}."
      },
      "description": "In-app notification template for order status updates",
      "required_variables": ["order_id", "status", "item_count", "total_amount", "status_message", "action_button"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787203796Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440006",
      "name": "Payment Reminder Template",
      "type": "in_app",
      "version": 1,
      "content": {
        "title": "Payment Reminder - ${{amount}}",
        "body": "Your payment of ${{amount}} is due on {{due_date}}.\n\nInvoice ID: {{invoice_id}}\n\nPlease complete your payment to avoid any service interruptions."
      },
      "description": "In-app notification template for payment reminders",
      "required_variables": ["amount", "due_date", "invoice_id"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787203879Z"
    }
  ]
```

## Best Practices

1. **Email Notifications**: Always include the `from` field with a valid email address
2. **Scheduled Notifications**: Use ISO 8601 format for `scheduled_at` timestamps
3. **Template Variables**: Ensure all required template variables are provided