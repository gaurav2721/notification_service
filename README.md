This service is used to send email , slack and in-app(for ios(apple push notification) and android(firebase cloud messaging)) notifications .

### Quick Start

### Build and Run with Docker
1. `make docker-build`
2. `make docker-run`

### View Docker Container Output
1. `make docker-exec`

## Building and Running

For detailed instructions on building and running the notification service, please refer to [Build.md](BUILD.md).

## Assumptions for this service

1. When a customer raises a notification request, only user IDs will be provided as recipients. The service will retrieve other necessary details from the pre-stored user information.
2. Only text-based content will be supported for notifications in this iteration.
3. Each notification request raised by customer/user will be linked to only one notification type/channel in this iteration.
4. In-app notifications will be limited to mobile push notifications for iOS and Android in this iteration.
5. All the information is stored in memory , database persistence may be added later
6. User and UserDeviceInfo have been preloaded and the apis for these are disabled by default , since it is considered to be out of scope
7. Interfaces for sending emails, Slack messages, and in-app notifications will be mocked in the first iteration to focus on building scalable service logic with features such as templates and scheduling.(If the .env does not have creds for the email, slack, apns and fcm , the information will be printed in a simple output/<service_name>.txt for eg output/email.txt)

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

Request Body (Immediate Notification without template):
```json
{
  "type": "email",
  "recipients": ["user-123", "user-456"],
  "content": {
    "subject": "Welcome To Tuskira",
    "email_body": "Hi! John Doe we welcome you to tuskira"
  },
  "from": {
    "email": "noreply@company.com"
  }
}
```

Request Body (Scheduled Notification without template):
```json
{
  "type": "email",
  "recipients": ["user-123", "user-456"],
  "content": {
    "subject": "Welcome To Tuskira",
    "email_body": "Hi! John Doe we welcome you to tuskira"
  },
  "scheduled_at": "2024-01-01T12:00:00Z",
  "from": {
    "email": "noreply@company.com"
  }
}
```

Request Body (Using Template):
```json
{
  "type": "email",
  "recipients": ["user-123", "user-456"],
  "template": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "data": {
      "name": "John Doe",
      "platform": "Tuskira"
    }
  },
  "from": {
    "email": "noreply@company.com"
  }
}
```

Request Body (Scheduled Notification using Template):
```json
{
  "type": "email",
  "recipients": ["user-123", "user-456"],
  "template": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "data": {
      "name": "John Doe",
      "platform": "Tuskira"
    }
  },
  "scheduled_at": "2024-01-01T12:00:00Z",
  "from": {
    "email": "noreply@company.com"
  }
}
```

**Request Parameters:**
- `type` (required): Notification type (`email`, `slack`, `in_app`)
- `recipients` (required): Array of user IDs (e.g., `["user-123", "user-456"]`)
- `content` (required if no template): Content object with notification-specific fields:
  - For email: `subject`, `email_body`
  - For slack: `text`
  - For in-app: `title`, `body`
- `template` (required if no content): Template object with:
  - `id` (string): Template identifier in UUID format
  - `data` (object): Key-value pairs for template parameters
- `scheduled_at` (optional): ISO 8601 timestamp (UTC) for scheduled delivery
- `from` (required for email notifications): Object containing sender information:
  - `email` (string): Sender email address

**Note**: Either `content` OR `template` must be provided, but not both. When using a template, the `data` object contains the key-value pairs that will replace the template variables.

**Important**: Email notifications require a mandatory `from` field with an email address. Non-email notifications should not include the `from` field.


## Notification Scheduling

The notification service supports two delivery modes:

### Immediate Delivery
- **Omit** the `scheduled_at` field or set it to `null`
- Notification will be sent immediately when the request is processed
- Best for real-time notifications, alerts, and instant communications

### Scheduled Delivery
- **Include** the `scheduled_at` field with an ISO 8601 timestamp (UTC)
- Notification will be stored and delivered at the specified time
- Best for reminders, scheduled announcements, and time-sensitive campaigns

**Note**: The `scheduled_at` timestamp should be an ISO 8601 timestamp (UTC) in the future. Past timestamps will result in immediate delivery.

**Important**: Recipients are always specified as user IDs (e.g., `["user-123", "user-456"]`).

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
    "scheduled_at": "2024-01-01T12:00:00Z",
    "from": {
      "email": "noreply@company.com"
    }
  }'
```

### Templates

#### Create Template
```
POST /api/v1/templates
```

Creates a new template with a higher logical version number.

**Request Body:**
```json
{
  "name": "Welcome Email Template",
  "type": "email",
  "content": {
    "subject": "Welcome to {{platform}}, {{name}}!",
    "email_body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
  },
  "required_variables": ["name", "platform", "username", "email", "account_type", "activation_link"],
  "description": "Welcome email template for new user onboarding"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Welcome Email Template",
  "type": "email",
  "version": 5,
  "status": "created",
  "created_at": "2024-01-01T10:00:00Z"
}
```

**Request Parameters:**
- `name` (string, required): Human-readable name for the template
- `type` (string, required): Notification type - must be one of: `email`, `slack`, `in_app`
- `content` (object, required): Template content with type-specific fields:
  - For email: `subject` (string), `email_body` (string)
  - For slack: `text` (string)
  - For in-app: `title` (string), `body` (string)
- `required_variables` (array, required): List of variable names that must be provided when using this template
- `description` (string, optional): Human-readable description of the template's purpose

**Response Parameters:**
- `id` (string): Unique template identifier in UUID format
- `name` (string): Template name as provided in request
- `type` (string): Notification type as provided in request
- `version` (integer): Logical version number assigned to this template
- `status` (string): Template status - "created" for new templates
- `created_at` (string): ISO 8601 timestamp when template was created

#### Get Template
```
GET /api/v1/templates/{templateId}/versions/{version}
```

Gets a specific version of a template.

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Welcome Email Template",
  "type": "email",
  "version": 1,
  "content": {
    "subject": "Welcome to {{platform}}, {{name}}!",
    "email_body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
  },
  "required_variables": ["name", "platform", "username", "email", "account_type", "activation_link"],
  "description": "Welcome email template for new user onboarding",
  "status": "active",
  "created_at": "2024-01-01T10:00:00Z"
}
```

**Path Parameters:**
- `templateId` (string, required): Unique identifier of the template
- `version` (integer, required): Specific version number to retrieve

**Response Parameters:**
- `id` (string): Template identifier in UUID format
- `name` (string): Human-readable template name
- `type` (string): Notification type (email, slack, in_app)
- `version` (integer): Template version number
- `content` (object): Template content with type-specific fields:
  - For email: `subject` (string), `email_body` (string)
  - For slack: `text` (string)
  - For in-app: `title` (string), `body` (string)
- `required_variables` (array): List of required variable names
- `description` (string): Template description
- `status` (string): Template status - "active" for available templates
- `created_at` (string): ISO 8601 timestamp when template was created

**Template Parameters:**
- `name` (required): Template name for identification
- `type` (required): Notification type (`email`, `slack`, `in_app`)
- `content` (required): Content object with notification-specific fields:
  - For email: `subject`, `email_body`
  - For slack: `text`
  - For in-app: `title`, `body`
- `required_variables` (required): Array of required variable names that must be provided
- `description` (optional): Template description for documentation

**Template Features:**
- **Immutable Versions**: Each template version is immutable and cannot be modified
- **Variable Validation**: Required variables are validated before sending notifications
- **Stable Payloads**: Template structure remains stable across versions
- **Content Encapsulation**: Content is properly encapsulated in the `content` key

## Template Management

### Immutable Versioning
- Each template has a version number as integer (e.g., 1, 2, 3)
- Template versions are immutable and cannot be modified once created
- New versions are created by incrementing the version number
- Old versions remain available for backward compatibility

### Variable Validation
- `required_variables`: Must be provided when using the template
- Validation occurs before notification sending
- Missing required variables result in an error response

### Template Lifecycle
1. **Create**: Register a new template using POST (automatically assigns next logical version)
2. **Retrieve**: Get specific template version using template ID and version number
3. **Version Management**: Each version is immutable once created

### API Behavior
- **Template Creation**: POST creates a new template with next available version number
- **Version-Specific Access**: GET requires both template ID and version number
- **Immutable Versions**: Template versions cannot be modified once created
- **Simple Structure**: Only two endpoints for complete template management

### Content Structure
- Templates use the same `content` structure as notifications
- Content is properly encapsulated and type-specific
- Supports all notification types: email, slack, in_app

## Template Examples

### 1. Create Email Templates

#### Welcome Email Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Welcome Email Template",
    "type": "email",
    "content": {
      "subject": "Welcome to {{platform}}, {{name}}!",
      "email_body": "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
    },
    "required_variables": ["name", "platform", "username", "email", "account_type", "activation_link"],
    "description": "Welcome email template for new user onboarding"
  }'
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Welcome Email Template",
  "type": "email",
  "version": 1,
  "status": "created",
  "created_at": "2024-01-01T10:00:00Z"
}
```

#### Password Reset Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Password Reset Template",
    "type": "email",
    "content": {
      "subject": "Password Reset Request - {{platform}}",
      "email_body": "Hello {{name}},\n\nWe received a request to reset your password for your {{platform}} account.\n\nIf you made this request, please click the link below to reset your password:\n{{reset_link}}\n\nThis link will expire in {{expiry_hours}} hours.\n\nIf you did not request a password reset, please ignore this email.\n\nBest regards,\nThe {{platform}} Security Team"
    },
    "required_variables": ["name", "platform", "reset_link", "expiry_hours"],
    "description": "Password reset email template"
  }'
```

#### Order Confirmation Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Confirmation Template",
    "type": "email",
    "content": {
      "subject": "Order Confirmed - #{{order_id}}",
      "email_body": "Dear {{customer_name}},\n\nThank you for your order! Your order has been confirmed and is being processed.\n\n**Order Details:**\n- Order ID: #{{order_id}}\n- Order Date: {{order_date}}\n- Total Amount: ${{total_amount}}\n- Payment Method: {{payment_method}}\n\n**Items Ordered:**\n{{items_list}}\n\n**Shipping Information:**\n{{shipping_address}}\n\n**Estimated Delivery:** {{delivery_date}}\n\nTrack your order: {{tracking_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team"
    },
    "required_variables": ["customer_name", "order_id", "order_date", "total_amount", "payment_method", "items_list", "shipping_address", "delivery_date", "tracking_link", "platform"],
    "description": "Order confirmation email template"
  }'
```

### 2. Create Slack Templates

#### System Alert Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "System Alert Template",
    "type": "slack",
    "content": {
      "text": "🚨 *{{alert_type}} Alert*\n\n*System:* {{system_name}}\n*Severity:* {{severity}}\n*Environment:* {{environment}}\n*Message:* {{message}}\n*Timestamp:* {{timestamp}}\n*Action Required:* {{action_required}}\n\n*Affected Services:* {{affected_services}}\n*Dashboard:* {{dashboard_link}}\n\nPlease take immediate action if this is a critical alert."
    },
    "required_variables": ["alert_type", "system_name", "severity", "environment", "message", "timestamp", "action_required", "affected_services", "dashboard_link"],
    "description": "Slack alert template for system monitoring"
  }'
```

#### Deployment Notification Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deployment Notification Template",
    "type": "slack",
    "content": {
      "text": "🚀 *Deployment {{status}}*\n\n*Service:* {{service_name}}\n*Environment:* {{environment}}\n*Version:* {{version}}\n*Deployed By:* {{deployed_by}}\n*Duration:* {{duration}}\n\n*Changes:*\n{{changes_summary}}\n\n*Rollback:* {{rollback_command}}\n*Monitoring:* {{monitoring_link}}"
    },
    "required_variables": ["status", "service_name", "environment", "version", "deployed_by", "duration", "changes_summary", "rollback_command", "monitoring_link"],
    "description": "Slack notification template for deployment events"
  }'
```

### 3. Create In-App Templates

#### Order Status Update Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Order Status Update Template",
    "type": "in_app",
    "content": {
      "title": "Order #{{order_id}} - {{status}}",
      "body": "Your order has been {{status}}.\n\n*Order Details:*\n- Items: {{item_count}} items\n- Total: ${{total_amount}}\n- Status: {{status}}\n\n{{status_message}}\n\n{{action_button}}"
    },
    "required_variables": ["order_id", "status", "item_count", "total_amount", "status_message", "action_button"],
    "description": "In-app notification template for order status updates"
  }'
```

#### Payment Reminder Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Payment Reminder Template",
    "type": "in_app",
    "content": {
      "title": "Payment Due - ${{amount}}",
      "body": "Your payment of ${{amount}} is due on {{due_date}}.\n\n*Invoice Details:*\n- Invoice #: {{invoice_id}}\n- Due Date: {{due_date}}\n- Amount: ${{amount}}\n\nPlease update your payment method or contact support if you have any questions."
    },
    "required_variables": ["amount", "due_date", "invoice_id"],
    "description": "In-app notification template for payment reminders"
  }'
```

### 4. Send Notifications Using Templates

#### Send Welcome Email Using Template
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
    "recipients": ["user-123", "user-456"]
  }'
```

#### Send System Alert Using Template
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
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
    "recipients": ["user-789", "user-101"]
  }'
```

#### Send Order Update Using Template
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "data": {
        "order_id": "ORD-2024-001",
        "status": "Shipped",
        "item_count": "3",
        "total_amount": "299.99",
        "status_message": "Your order has been shipped and is on its way!",
        "action_button": "Track Order"
      }
    },
    "recipients": ["user-123"]
  }'
```

### 5. Schedule Template-Based Notifications

#### Schedule Welcome Email
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
    "recipients": ["user-789"],
    "scheduled_at": "2024-01-01T12:00:00Z"
  }'
```

#### Schedule Payment Reminder
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "data": {
        "amount": "99.99",
        "due_date": "2024-01-15",
        "invoice_id": "INV-2024-001"
      }
    },
    "recipients": ["user-456"],
    "scheduled_at": "2024-01-15T12:00:00Z"
  }'
```

### 6. Template Version Management

#### Create New Version of Welcome Template
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Welcome Email Template",
    "type": "email",
    "content": {
      "subject": "Welcome to {{platform}}, {{name}}! 🎉",
      "email_body": "Hello {{name}},\n\n🎉 Welcome to {{platform}}! We are thrilled to have you join our community.\n\nYour account has been successfully created:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\n🔗 Activate your account: {{activation_link}}\n\n💡 Get started with our quick guide: {{getting_started_link}}\n\nIf you have any questions, our support team is here to help!\n\nBest regards,\nThe {{platform}} Team"
    },
    "required_variables": ["name", "platform", "username", "email", "account_type", "activation_link", "getting_started_link"],
    "description": "Updated welcome email template with enhanced styling and additional resources"
  }'
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Welcome Email Template",
  "type": "email",
  "version": 2,
  "status": "created",
  "created_at": "2024-01-01T11:00:00Z"
}
```

### 7. Retrieve Template Versions

#### Get Specific Template Version
```bash
curl -X GET http://localhost:8080/api/v1/templates/550e8400-e29b-41d4-a716-446655440000/versions/1
```

#### Get Latest Template Version
```bash
curl -X GET http://localhost:8080/api/v1/templates/550e8400-e29b-41d4-a716-446655440000/versions/2
```

## Notification Types

- `email`: Send notifications via email using SMTP
- `slack`: Send notifications to Slack channels
- `in_app`: Store notifications for in-app display

## Building and Running

For detailed instructions on building and running the notification service, please refer to [Build.md](Build.md).

The Build.md file contains:
- Docker build and run commands
- Local development setup
- Environment variable configuration
- Testing instructions
- Debug mode setup

### Quick Start

1. **Clone the repository**
2. **Set up environment variables** (see [Build.md](Build.md) for detailed configuration)
3. **Build and run** using one of the following methods:
   - **Docker**: `make docker-build && make docker-run`
   - **Local**: `make build && make run`
   - **Debug**: `make run-debug`

## Development

```
notification_service/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── Makefile                # Build and deployment scripts
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── env.example             # Environment variables template
├── README.md               # Project documentation
├── USER_README.md          # User-focused documentation
├── bin/                    # Compiled binaries
│   └── notification-service # Executable binary
├── docs/                   # Documentation files
├── examples/               # Example implementations
│   ├── user_service_integration.go      # User service integration example
│   └── service_architecture_example.go  # Service architecture example
├── models/                 # Data models and types
│   └── notification.go     # Notification data models
├── handlers/               # HTTP request handlers
│   └── notification_handlers.go  # Notification API handlers
├── routes/                 # API routing configuration
│   ├── routes.go           # Main routing setup
│   ├── notification_routes.go    # Notification endpoints
│   ├── user_routes.go      # User management endpoints
│   ├── template_routes.go  # Template management endpoints
│   ├── health_routes.go    # Health check endpoints
│   └── middleware/         # HTTP middleware
│       └── middleware.go   # Authentication and logging middleware
├── services/               # Core service implementations
│   ├── interfaces.go       # Service interfaces
│   ├── notification_service.go  # Main notification service
│   ├── email_service.go    # Email service implementation
│   ├── slack_service.go    # Slack service implementation
│   ├── inapp_service.go    # In-app service implementation
│   └── scheduler_service.go # Scheduler service implementation
├── notification_manager/   # Notification orchestration
│   ├── interface.go        # Manager interface
│   ├── notification_manager.go    # Main manager implementation
│   ├── notification_manager_test.go  # Manager tests
│   ├── config.go           # Manager configuration
│   └── errors.go           # Manager error definitions
├── external_services/      # External service integrations
│   ├── service_factory.go  # Service factory pattern
│   ├── service_provider.go # Service provider implementation
│   ├── email/              # Email service module
│   │   ├── interface.go    # Email service interface
│   │   ├── email_service.go # Email service implementation
│   │   ├── config.go       # Email configuration
│   │   └── errors.go       # Email error definitions
│   ├── slack/              # Slack service module
│   │   ├── interface.go    # Slack service interface
│   │   ├── slack_service.go # Slack service implementation
│   │   ├── config.go       # Slack configuration
│   │   └── errors.go       # Slack error definitions
│   ├── inapp/              # In-app notification module
│   │   ├── interface.go    # In-app service interface
│   │   ├── inapp_service.go # In-app service implementation
│   │   ├── config.go       # In-app configuration
│   │   └── errors.go       # In-app error definitions
│   └── user/               # User service module
│       ├── interface.go    # User service interface
│       ├── user_service.go # User service implementation
│       ├── user_service_test.go # User service tests
│       ├── config.go       # User service configuration
│       └── errors.go       # User service error definitions
└── helper/                 # Utility and helper functions
    └── scheduler/          # Scheduler utilities
        ├── interface.go    # Scheduler interface
        ├── scheduler.go    # Scheduler implementation
        ├── config.go       # Scheduler configuration
        └── errors.go       # Scheduler error definitions
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