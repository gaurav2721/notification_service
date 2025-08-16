# Notification Service - Run Instructions

This document provides step-by-step instructions to run the notification service and test all its features.

## Prerequisites

- Docker and Docker Compose installed
- curl (for API testing)

## 1. Build and Run the Service

### 1.1 Build Docker Image
```bash
make docker-build
```

### 1.2 Run Docker Container
```bash
make docker-run
```

## 2. Get All Users

The service comes with some pre-loaded users for testing purposes.

**Note:** The service redirects `/api/v1/users` to `/api/v1/users/` (with trailing slash), so use the trailing slash version.

```bash
curl -X GET http://localhost:8080/api/v1/users/
```

**Expected Output:**
```json
{
  "count": 8,
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
}
```

## 3. Immediate Notifications

### 3.1 Send Immediate Email Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
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

**Expected Output:**
```json
{
  "id": "1755282275013750680",
  "status": "sent",
  "message": "Email notification sent successfully",
  "sent_at": "2025-08-15T18:24:35.013754596Z",
  "channel": "email"
}
```

### 3.2 Send Immediate Slack Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "content": {
      "text": "*System Alert*\n\n*Service:* Notification Service\n*Status:* Running\n*Environment:* Development\n*Message:* All systems operational"
    },
    "recipients": ["user-001"]
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901235",
  "status": "sent",
  "message": "Slack notification sent successfully",
  "sent_at": "2024-01-15T12:01:00Z",
  "channel": "slack"
}
```

### 3.3 Send Immediate In-App Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "content": {
      "title": "New Feature Available",
      "body": "We have just released a new feature! Check it out in your dashboard."
    },
    "recipients": ["user-001"]
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901236",
  "status": "sent",
  "message": "In-app notification sent successfully",
  "sent_at": "2024-01-15T12:02:00Z",
  "channel": "in_app"
}
```

## 4. Scheduled Notifications

### 4.1 Schedule Email Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Reminder: Complete Your Profile",
      "email_body": "Hello,\n\nThis is a friendly reminder to complete your profile information.\n\nBest regards,\nThe Team"
    },
    "recipients": ["user-001"],
    "scheduled_at": "2024-01-15T14:00:00Z",
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901237",
  "status": "scheduled",
  "message": "Email notification scheduled successfully",
  "scheduled_at": "2024-01-15T14:00:00Z",
  "channel": "email"
}
```

### 4.2 Schedule Slack Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "content": {
      "text": "📅 *Daily Standup Reminder*\n\nTime: 9:00 AM\nChannel: #daily-standup\nAgenda: Project updates and blockers"
    },
    "recipients": ["user-001", "user-002", "user-003"],
    "scheduled_at": "2024-01-16T09:00:00Z"
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901238",
  "status": "scheduled",
  "message": "Slack notification scheduled successfully",
  "scheduled_at": "2024-01-16T09:00:00Z",
  "channel": "slack"
}
```

### 4.3 Schedule In-App Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "content": {
      "title": "Weekly Report Ready",
      "body": "Your weekly performance report is now available. Click here to view it."
    },
    "recipients": ["user-001"],
    "scheduled_at": "2024-01-20T08:00:00Z"
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901239",
  "status": "scheduled",
  "message": "In-app notification scheduled successfully",
  "scheduled_at": "2024-01-20T08:00:00Z",
  "channel": "in_app"
}
```

## 5. Get Predefined Templates

The service comes with several predefined templates for common use cases.

```bash
curl -X GET http://localhost:8080/api/v1/templates/predefined
```

**Expected Output:**
```json
{
  "count": 7,
  "templates": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Welcome Email Template",
      "type": "email",
      "version": 1,
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
      "description": "In-app notification template for payment reminders",
      "required_variables": ["amount", "due_date", "invoice_id"],
      "status": "active",
      "created_at": "2025-08-15T18:23:46.787203879Z"
    }
  ]
}
```

## 6. Use Predefined Templates for Immediate Notifications

### 6.1 Send Immediate Email Using Welcome Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "data": {
        "name": "John Doe",
        "platform": "Tuskira",
        "username": "johndoe",
        "email": "john.doe@example.com",
        "account_type": "Premium",
        "activation_link": "https://tuskira.com/activate?token=abc123def456"
      }
    },
    "recipients": ["user-001"],
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901240",
  "status": "sent",
  "message": "Email notification sent successfully using template",
  "sent_at": "2024-01-15T12:05:00Z",
  "channel": "email"
}
```

### 6.2 Send Immediate Slack Using System Alert Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "data": {
        "alert_type": "Database Connection",
        "system_name": "User Service",
        "severity": "Critical",
        "environment": "Production",
        "message": "Database connection timeout after 30 seconds",
        "timestamp": "2024-01-15T12:06:00Z",
        "action_required": "Immediate investigation required",
        "affected_services": "User authentication, Profile management",
        "dashboard_link": "https://dashboard.example.com/alerts"
      }
    },
    "recipients": ["user-001", "user-002"]
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901241",
  "status": "sent",
  "message": "Slack notification sent successfully using template",
  "sent_at": "2024-01-15T12:06:00Z",
  "channel": "slack"
}
```

### 6.3 Send Immediate In-App Using Order Status Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "data": {
        "order_id": "ORD-2024-001",
        "status": "Shipped",
        "item_count": 3,
        "total_amount": "299.99",
        "status_message": "Your order has been shipped and is on its way!",
        "action_button": "Track Order"
      }
    },
    "recipients": ["user-001"]
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901242",
  "status": "sent",
  "message": "In-app notification sent successfully using template",
  "sent_at": "2024-01-15T12:07:00Z",
  "channel": "in_app"
}
```

## 7. Use Predefined Templates for Scheduled Notifications

### 7.1 Schedule Email Using Welcome Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "data": {
        "name": "Jane Smith",
        "platform": "Tuskira",
        "username": "janesmith",
        "email": "jane.smith@example.com",
        "account_type": "Standard",
        "activation_link": "https://tuskira.com/activate?token=def456ghi789"
      }
    },
    "recipients": ["user-456"],
    "scheduled_at": "2024-01-15T15:00:00Z",
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901243",
  "status": "scheduled",
  "message": "Email notification scheduled successfully using template",
  "scheduled_at": "2024-01-15T15:00:00Z",
  "channel": "email"
}
```

### 7.2 Schedule Slack Using System Alert Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "data": {
        "alert_type": "Backup Status",
        "system_name": "Database Backup",
        "severity": "Info",
        "environment": "Production",
        "message": "Daily backup completed successfully",
        "timestamp": "2024-01-16T02:00:00Z",
        "action_required": "No action required",
        "affected_services": "All services",
        "dashboard_link": "https://dashboard.example.com/backups"
      }
    },
    "recipients": ["user-001", "user-002", "user-003"],
    "scheduled_at": "2024-01-16T02:30:00Z"
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901244",
  "status": "scheduled",
  "message": "Slack notification scheduled successfully using template",
  "scheduled_at": "2024-01-16T02:30:00Z",
  "channel": "slack"
}
```

### 7.3 Schedule In-App Using Order Status Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "data": {
        "order_id": "ORD-2024-002",
        "status": "Delivered",
        "item_count": 1,
        "total_amount": "99.99",
        "status_message": "Your order has been delivered! Please rate your experience.",
        "action_button": "Rate Order"
      }
    },
    "recipients": ["user-002"],
    "scheduled_at": "2024-01-16T10:00:00Z"
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901245",
  "status": "scheduled",
  "message": "In-app notification scheduled successfully using template",
  "scheduled_at": "2024-01-16T10:00:00Z",
  "channel": "in_app"
}
```

## 8. Define a New Email Template

Create a custom email template for password reset notifications.

```bash
curl -X POST http://localhost:8080/api/v1/templates \
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

**Expected Output:**
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

## 9. Use the New Template for Immediate and Scheduled Notifications

### 9.1 Send Immediate Email Using New Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "template-password-reset-custom",
      "data": {
        "user_name": "John Doe",
        "platform_name": "Tuskira",
        "reset_link": "https://tuskira.com/reset-password?token=xyz789abc123",
        "expiry_hours": 24
      }
    },
    "recipients": ["user-123"],
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901246",
  "status": "sent",
  "message": "Email notification sent successfully using template",
  "sent_at": "2024-01-15T12:11:00Z",
  "channel": "email"
}
```

### 9.2 Schedule Email Using New Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "template-password-reset-custom",
      "data": {
        "user_name": "Jane Smith",
        "platform_name": "Tuskira",
        "reset_link": "https://tuskira.com/reset-password?token=def456ghi789",
        "expiry_hours": 24
      }
    },
    "recipients": ["user-002"],
    "scheduled_at": "2024-01-15T16:00:00Z",
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

**Expected Output:**
```json
{
  "id": "1705312345678901247",
  "status": "scheduled",
  "message": "Email notification scheduled successfully using template",
  "scheduled_at": "2024-01-15T16:00:00Z",
  "channel": "email"
}
```

## 10. Check Notification Status

You can check the status of any notification using its ID.

```bash
curl -X GET http://localhost:8080/api/v1/notifications/1705312345678901234
```

**Expected Output:**
```json
{
  "id": "1705312345678901234",
  "status": "sent",
  "message": "Email notification sent successfully",
  "sent_at": "2024-01-15T12:00:00Z",
  "channel": "email",
  "recipients": ["user-001", "user-002"]
}
```

## 11. Health Check

Check if the service is running properly.

```bash
curl -X GET http://localhost:8080/health
```

**Expected Output:**
```json
{
  "service": "notification-service",
  "status": "healthy",
  "timestamp": "2025-08-15T18:23:52.426265799Z"
}
```

## Additional Notes

- All timestamps are in ISO 8601 format (UTC)
- User IDs are used as recipients instead of email addresses for better security
- The service supports three notification types: `email`, `slack`, and `in_app`
- **Email notifications require a mandatory "from" field with an email address**
- **Non-email notifications should not include the "from" field**
- Templates support variable substitution using `{{variable_name}}` syntax
- Scheduled notifications are processed by the service's internal scheduler
- All API responses include appropriate HTTP status codes
- Error responses include detailed error messages for debugging

## Troubleshooting

If you encounter issues:

1. **Service not starting**: Check if port 8080 is available
2. **Template not found**: Verify the template ID exists using the predefined templates endpoint
3. **User not found**: Ensure the user ID exists in the system
4. **Scheduled notifications not sending**: Check the service logs for scheduler errors

## Environment Configuration

Make sure to set up your environment variables in a `.env` file based on `env.example`:

```bash
cp env.example .env
# Edit .env with your actual configuration values
```
