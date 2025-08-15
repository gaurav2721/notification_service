## Installation

### Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd notification_service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your configuration
```

4. Run the service:
```bash
go run main.go
```

### Using Docker

1. Build and run with Docker Compose:
```bash
docker-compose up --build
```

2. Or build and run manually:
```bash
docker build -t notification-service .
docker run -p 8080:8080 --env-file .env notification-service
```

## Configuration

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Slack Configuration
SLACK_BOT_TOKEN=xoxb-your-slack-bot-token
SLACK_CHANNEL_ID=C1234567890
```

## API Endpoints

### Health Check
```
GET /health
```

**Description:**
The health check endpoint provides real-time status information about the notification service. It verifies that the service is running and responsive, confirming that the notification service itself is operational.

**cURL Command:**
```bash
curl -X GET http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T10:00:00Z",
  "service": "notification-service"
}
```

**Response Fields:**
- `status` (string): Overall health status (`healthy`, `unhealthy`, `degraded`)
- `timestamp` (string): ISO 8601 timestamp of the health check
- `service` (string): Service name identifier

**HTTP Status Codes:**
- `200 OK`: Service is healthy and operational
- `503 Service Unavailable`: Service is unhealthy or not responding
- `500 Internal Server Error`: Health check itself failed

**Use Cases:**
- Load balancer health checks
- Monitoring and alerting systems
- Container orchestration (Kubernetes, Docker Swarm)
- Service mesh health verification

### Notifications

#### Send Notification
```
POST /api/v1/notifications
```

Request Body (Immediate Notification):
```json
{
  "type": "email",
  "recipients": ["user-123", "user-456"],
  "content": {
    "subject": "Welcome To Tuskira",
    "email_body": "Hi! John Doe we welcome you to tuskira"
  }
}
```

Request Body (Scheduled Notification):
```json
{
  "type": "email",
  "recipients": ["user-123", "user-456"],
  "content": {
    "subject": "Welcome To Tuskira",
    "email_body": "Hi! John Doe we welcome you to tuskira"
  },
  "scheduled_at": "2024-01-01T10:00:00Z"
}
```

**Request Parameters:**
- `type` (required): Notification type (`email`, `slack`, `in_app`)
- `recipients` (required): Array of user IDs (e.g., `["user-123", "user-456"]`)
- `content` (required): Content object with notification-specific fields:
  - For email: `subject`, `email_body`
  - For slack: `text`
  - For in-app: `title`, `body`
- `template` (optional): Template object with `id` and `data` for template-based notifications
- `metadata` (optional): Additional data as key-value pairs
- `scheduled_at` (optional): ISO 8601 timestamp for scheduled delivery

#### Get Notification Status
```
GET /api/v1/notifications/{id}
```

### Templates

#### Create Template
```
POST /api/v1/templates
```

Request Body:
```json
{
  "name": "Welcome Email Template",
  "type": "email",
  "subject": "Welcome to {{platform}}, {{name}}!",
  "body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team",
  "variables": ["name", "platform", "username", "email", "account_type", "activation_link"]
}
```

**Template Parameters:**
- `name` (required): Template name for identification
- `type` (required): Notification type (`email`, `slack`, `in_app`)
- `subject` (required): Template subject/title with variable placeholders
- `body` (required): Template body with variable placeholders using `{{variable_name}}` syntax
- `variables` (optional): Array of variable names used in the template

**Variable Syntax:**
- Use `{{variable_name}}` in subject and body for dynamic content
- Variables are replaced with actual values when sending notifications
- Common variables: `{{name}}`, `{{email}}`, `{{platform}}`, `{{order_id}}`, etc.

#### Get Template
```
GET /api/v1/templates/{id}
```

#### Update Template
```
PUT /api/v1/templates/{id}
```

#### Delete Template
```
DELETE /api/v1/templates/{id}
```

## Notification Scheduling

The notification service supports two delivery modes:

### Immediate Delivery
- **Omit** the `scheduled_at` field or set it to `null`
- Notification will be sent immediately when the request is processed
- Best for real-time notifications, alerts, and instant communications

### Scheduled Delivery
- **Include** the `scheduled_at` field with an ISO 8601 timestamp
- Notification will be queued and delivered at the specified time
- Best for reminders, scheduled announcements, and time-sensitive campaigns

**Note**: The `scheduled_at` timestamp should be in the future. Past timestamps will result in immediate delivery.

**Important**: Recipients are always specified as user IDs (e.g., `["user-123", "user-456"]`). The service will automatically resolve the appropriate delivery channels (email, Slack, in-app) for each user based on their preferences and settings.

## Usage Examples

### Send an Email Notification (Immediate)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Test Email",
      "email_body": "This is a test email notification"
    },
    "recipients": ["user-123", "user-456"]
  }'
```

### Send a Slack Notification (Immediate)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "content": {
      "text": "This is a test Slack notification"
    },
    "recipients": ["user-123", "user-456"]
  }'
```

### Send an In-App Notification (Immediate)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "content": {
      "title": "New Message",
      "body": "You have received a new message from John"
    },
    "recipients": ["user-123", "user-456"]
  }'
```

### Schedule a Notification (Future Delivery)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Scheduled Email",
      "email_body": "This email was scheduled for future delivery"
    },
    "recipients": ["user-123", "user-456"],
    "scheduled_at": "2024-01-01T10:00:00Z"
  }'
```

## Template Examples

### Create Email Template

```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Welcome Email Template",
    "type": "email",
    "subject": "Welcome to {{platform}}, {{name}}!",
    "body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team",
    "variables": ["name", "platform", "username", "email", "account_type", "activation_link"]
  }'
```

### Create Slack Template

```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alert Notification Template",
    "type": "slack",
    "subject": "",
    "body": "🚨 *{{alert_type}} Alert*\n\n*System:* {{system_name}}\n*Severity:* {{severity}}\n*Message:* {{message}}\n*Timestamp:* {{timestamp}}\n*Action Required:* {{action_required}}\n\nPlease take immediate action if this is a critical alert.",
    "variables": ["alert_type", "system_name", "severity", "message", "timestamp", "action_required"]
  }'
```

### Create In-App Template

```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Confirmation Template",
    "type": "in_app",
    "subject": "",
    "body": "🎉 *Order Confirmed!*\n\nYour order #{{order_id}} has been successfully placed.\n\n*Order Details:*\n- Items: {{item_count}} items\n- Total Amount: ${{total_amount}}\n- Estimated Delivery: {{delivery_date}}\n\nTrack your order: {{tracking_link}}",
    "variables": ["order_id", "item_count", "total_amount", "delivery_date", "tracking_link"]
  }'
```

### Send Notification Using Template (Email)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "welcome_email_template",
      "data": {
        "name": "John Doe",
        "platform": "Acme Corp",
        "username": "johndoe",
        "email": "john.doe@example.com",
        "account_type": "Premium",
        "activation_link": "https://acme.com/activate?token=abc123"
      }
    },
    "recipients": ["user-123", "user-456"]
  }'
```

### Send Notification Using Template (Slack)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "template": {
      "id": "alert_notification_template",
      "data": {
        "alert_type": "System Down",
        "system_name": "Payment Gateway",
        "severity": "Critical",
        "message": "Payment processing is currently unavailable",
        "timestamp": "2024-01-01T10:00:00Z",
        "action_required": "Immediate investigation required"
      }
    },
    "recipients": ["user-123", "user-456"]
  }'
```

### Send Notification Using Template (In-App)

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "template": {
      "id": "order_confirmation_template",
      "data": {
        "order_id": "ORD-2024-001",
        "item_count": "3",
        "total_amount": "299.99",
        "delivery_date": "2024-01-05",
        "tracking_link": "https://acme.com/track/ORD-2024-001"
      }
    },
    "recipients": ["user-123"]
  }'
```

### Schedule Template-Based Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "welcome_email_template",
      "data": {
        "name": "Jane Smith",
        "platform": "Acme Corp",
        "username": "janesmith",
        "email": "jane.smith@example.com",
        "account_type": "Standard",
        "activation_link": "https://acme.com/activate?token=def456"
      }
    },
    "recipients": ["user-789"],
    "scheduled_at": "2024-01-01T09:00:00Z"
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
    "body": "Dear {{customer_name}}, your order #{{order_id}} has been confirmed.",
    "variables": ["order_id", "customer_name"]
  }'
```

## Notification Types

- `email`: Send notifications via email using SMTP
- `slack`: Send notifications to Slack channels
- `in_app`: Store notifications for in-app display

## Development

### Project Structure

```
notification_service/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── env.example             # Environment variables template
├── models/
│   └── notification.go     # Data models and types
├── services/
│   ├── interfaces.go       # Service interfaces
│   ├── notification_service.go  # Main notification service
│   ├── email_service.go    # Email service implementation
│   ├── slack_service.go    # Slack service implementation
│   ├── inapp_service.go    # In-app service implementation
│   └── scheduler_service.go # Scheduler service implementation
└── handlers/
    └── notification_handlers.go  # HTTP request handlers
```

### Adding New Notification Channels

1. Create a new service implementation in `services/`
2. Implement the appropriate interface
3. Add the new channel type to the notification manager
4. Update the routing logic in `SendNotification`

### Testing

Run tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request