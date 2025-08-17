This file shows how to quickly test this service


## 1. Immediate Notifications

### 1.1 Send Immediate Email Notification

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

**Expected Output:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "sent"
}
```

### 1.2 Send Immediate Slack Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
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
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "status": "sent"
}
```

### 1.3 Send Immediate In-App Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
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
  "id": "777e8901-e89b-12d3-a456-426614174019",
  "status": "sent"
}
```

## 2. Scheduled Notifications

### 2.1 Schedule Email Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
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
  "id": "789e0123-e89b-12d3-a456-426614174002",
  "status": "scheduled"
}
```

## 3. Use Predefined Templates for Immediate Notifications

### 3.1 Send Immediate Email Using Welcome Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "version": 1,
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
  "id": "101e2345-e89b-12d3-a456-426614174003",
  "status": "sent"
}
```

## 4. Use Predefined Templates for Scheduled Notifications

### 4.1 Schedule Email Using Welcome Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "version": 1,
      "data": {
        "name": "Jane Smith",
        "platform": "Tuskira",
        "username": "janesmith",
        "email": "jane.smith@example.com",
        "account_type": "Standard",
        "activation_link": "https://tuskira.com/activate?token=def456ghi789"
      }
    },
    "recipients": ["user-001"],
    "scheduled_at": "2024-01-15T15:00:00Z",
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

**Expected Output:**
```json
{
  "id": "202e3456-e89b-12d3-a456-426614174004",
  "status": "scheduled"
}
```

## 5. Define a New Email Template

Create a custom email template for password reset notifications.

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

**Expected Output:**
```json
{
  "id": "43138245-0467-49e4-a3cd-fef1d6b690f3",
  "name": "Password Reset Template",
  "type": "email",
  "version": 1,
  "status": "created",
  "created_at": "2025-08-15T18:25:00Z"
}
```

## 6. Use the New Template for Immediate and Scheduled Notifications

### 6.1 Send Immediate Email Using New Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "43138245-0467-49e4-a3cd-fef1d6b690f3",
      "version": 1,
      "data": {
        "user_name": "John Doe",
        "platform_name": "Tuskira",
        "reset_link": "https://tuskira.com/reset-password?token=xyz789abc123",
        "expiry_hours": 24
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
  "id": "303e4567-e89b-12d3-a456-426614174005",
  "status": "sent"
}
```

### 6.2 Schedule Email Using New Template

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer gaurav" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "43138245-0467-49e4-a3cd-fef1d6b690f3",
      "version": 1,
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
  "id": "404e5678-e89b-12d3-a456-426614174006",
  "status": "scheduled"
}
```
