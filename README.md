# Notification Service

A comprehensive notification service built in Go that supports multiple notification channels including email, Slack, and in-app notifications. The service provides scheduling capabilities, template management, and a RESTful API.

## Features

- **Multiple Notification Channels**: Email, Slack, and In-app notifications
- **Notification Scheduling**: Schedule notifications for future delivery
- **Template Management**: Create, update, and manage notification templates
- **RESTful API**: Clean HTTP API for all operations
- **Priority Levels**: Support for different notification priorities
- **Metadata Support**: Additional data can be attached to notifications
- **Health Checks**: Built-in health monitoring

## Architecture

The service follows a clean architecture pattern with:

- **Models**: Data structures and types
- **Services**: Business logic and external integrations
- **Handlers**: HTTP request handling
- **Main**: Application entry point and configuration

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional)
- SMTP server credentials (for email notifications)
- Slack Bot Token (for Slack notifications)

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

### Notifications

#### Send Notification
```
POST /api/v1/notifications
```

Request Body:
```json
{
  "type": "email",
  "priority": "normal",
  "title": "Welcome!",
  "message": "Welcome to our platform!",
  "recipients": ["user@example.com"],
  "metadata": {
    "user_id": "123",
    "campaign": "welcome"
  },
  "scheduled_at": "2024-01-01T10:00:00Z"
}
```

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
  "name": "Welcome Email",
  "type": "email",
  "subject": "Welcome to {{platform}}",
  "body": "Hello {{name}}, welcome to {{platform}}!",
  "variables": ["name", "platform"]
}
```

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

## Usage Examples

### Send an Email Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Test Email",
    "message": "This is a test email notification",
    "recipients": ["test@example.com"]
  }'
```

### Send a Slack Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "title": "Slack Alert",
    "message": "This is a test Slack notification",
    "recipients": ["general"]
  }'
```

### Schedule a Notification

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Scheduled Email",
    "message": "This email was scheduled",
    "recipients": ["user@example.com"],
    "scheduled_at": "2024-01-01T10:00:00Z"
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

## Priority Levels

- `low`: Low priority notifications
- `normal`: Standard priority (default)
- `high`: High priority notifications
- `urgent`: Urgent notifications

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