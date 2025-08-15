# Notification Service API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, the API does not require authentication. In a production environment, you should implement proper authentication and authorization.

## Endpoints

### Health Check

#### GET /health
Check the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T10:00:00Z",
  "service": "notification-service"
}
```

### Notifications

#### POST /api/v1/notifications
Send a notification through the specified channel.

**Request Body:**
```json
{
  "type": "email",
  "priority": "normal",
  "title": "Welcome!",
  "message": "Welcome to our platform!",
  "recipients": ["user@example.com"],
  "template_id": "optional-template-id",
  "metadata": {
    "user_id": "123",
    "campaign": "welcome"
  },
  "scheduled_at": "2024-01-01T10:00:00Z"
}
```

**Parameters:**
- `type` (required): Notification type (`email`, `slack`, `in_app`)
- `priority` (optional): Priority level (`low`, `normal`, `high`, `urgent`)
- `title` (required): Notification title
- `message` (required): Notification message
- `recipients` (required): Array of recipient identifiers
- `template_id` (optional): ID of a template to use
- `metadata` (optional): Additional data as key-value pairs
- `scheduled_at` (optional): ISO 8601 timestamp for scheduled delivery

**Response:**
```json
{
  "id": "notification-id",
  "status": "sent",
  "message": "Notification sent successfully",
  "sent_at": "2024-01-01T10:00:00Z",
  "channel": "email"
}
```

#### GET /api/v1/notifications/{id}
Get the status of a notification.

**Parameters:**
- `id` (path): Notification ID

**Response:**
```json
{
  "id": "notification-id",
  "status": "sent",
  "message": "Notification sent successfully",
  "sent_at": "2024-01-01T10:00:00Z",
  "channel": "email"
}
```

### Templates

#### POST /api/v1/templates
Create a new notification template.

**Request Body:**
```json
{
  "name": "Welcome Email",
  "type": "email",
  "subject": "Welcome to {{platform}}",
  "body": "Hello {{name}}, welcome to {{platform}}!",
  "variables": ["name", "platform"]
}
```

**Parameters:**
- `name` (required): Template name
- `type` (required): Notification type (`email`, `slack`, `in_app`)
- `subject` (required): Template subject/title
- `body` (required): Template body with variable placeholders
- `variables` (optional): Array of variable names used in the template

**Response:**
```json
{
  "id": "template-id",
  "name": "Welcome Email",
  "type": "email",
  "subject": "Welcome to {{platform}}",
  "body": "Hello {{name}}, welcome to {{platform}}!",
  "variables": ["name", "platform"],
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

#### GET /api/v1/templates/{id}
Get a notification template.

**Parameters:**
- `id` (path): Template ID

**Response:**
```json
{
  "id": "template-id",
  "name": "Welcome Email",
  "type": "email",
  "subject": "Welcome to {{platform}}",
  "body": "Hello {{name}}, welcome to {{platform}}!",
  "variables": ["name", "platform"],
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

#### PUT /api/v1/templates/{id}
Update a notification template.

**Parameters:**
- `id` (path): Template ID

**Request Body:**
```json
{
  "name": "Updated Welcome Email",
  "type": "email",
  "subject": "Welcome to {{platform}}",
  "body": "Hello {{name}}, welcome to our updated {{platform}}!",
  "variables": ["name", "platform"]
}
```

**Response:**
```json
{
  "id": "template-id",
  "name": "Updated Welcome Email",
  "type": "email",
  "subject": "Welcome to {{platform}}",
  "body": "Hello {{name}}, welcome to our updated {{platform}}!",
  "variables": ["name", "platform"],
  "updated_at": "2024-01-01T10:00:00Z"
}
```

#### DELETE /api/v1/templates/{id}
Delete a notification template.

**Parameters:**
- `id` (path): Template ID

**Response:**
```json
{
  "message": "Template deleted successfully"
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Validation error message"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error message"
}
```

## Examples

### Send an Email Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Welcome to Our Platform",
    "message": "Thank you for joining us!",
    "recipients": ["user@example.com"],
    "priority": "normal"
  }'
```

### Send a Slack Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "title": "System Alert",
    "message": "Server is running low on memory",
    "recipients": ["#alerts"],
    "priority": "high"
  }'
```

### Schedule a Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Reminder",
    "message": "Don\'t forget about the meeting tomorrow",
    "recipients": ["user@example.com"],
    "scheduled_at": "2024-01-02T09:00:00Z"
  }'
```

### Create a Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Confirmation",
    "type": "email",
    "subject": "Order #{{order_id}} Confirmed",
    "body": "Dear {{customer_name}}, your order #{{order_id}} has been confirmed and will be shipped soon.",
    "variables": ["order_id", "customer_name"]
  }'
```

## Rate Limiting

Currently, there are no rate limits implemented. In a production environment, you should implement appropriate rate limiting to prevent abuse.

## Versioning

The API is versioned using the URL path (`/api/v1/`). Future versions will be available at `/api/v2/`, etc.

## Support

For API support and questions, please refer to the project documentation or open an issue in the repository. 